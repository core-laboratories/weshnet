package weshnet

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	ma "github.com/multiformats/go-multiaddr"
)

// NetworkClassification describes what an observed connection path proves.
// It does not describe peer presence, reachability, or message delivery.
type NetworkClassification int

const (
	NetworkClassificationUnknown NetworkClassification = iota
	NetworkClassificationInternet
	NetworkClassificationOffline
)

type NetworkConnectionState int

const (
	NetworkConnectionStateUnknown NetworkConnectionState = iota
	NetworkConnectionStateConnected
	NetworkConnectionStateDisconnected
)

// NetworkConnectionEvidence contains facts supplied by an observer-side
// integration. The source must only inspect a connection; it must not alter it.
// Standard libp2p transport and address metadata alone does not set either flag.
type NetworkConnectionEvidence struct {
	InternetRequired    bool
	VerifiedLocalDirect bool
	Details             []string
}

// NetworkConnectionEvidenceSource can add independently verified path facts.
// Wesh does not contact an external service to obtain this evidence.
type NetworkConnectionEvidenceSource interface {
	EvidenceForConnection(network.Conn) NetworkConnectionEvidence
}

type NetworkConnectionMetadata struct {
	PeerID          string
	ConnectionID    string
	State           NetworkConnectionState
	LocalMultiaddr  string
	RemoteMultiaddr string
	Transport       string
	Direct          bool
	Relayed         bool
	ObservedAt      time.Time
	Evidence        NetworkConnectionEvidence
}

type NetworkConnectionStatus struct {
	NetworkConnectionMetadata
	Classification  NetworkClassification
	EvidenceDetails []string
}

type NetworkPeerStatus struct {
	PeerID          string
	Classification  NetworkClassification
	ConnectionIDs   []string
	Mixed           bool
	EvidenceDetails []string
}

type NetworkStatusSnapshot struct {
	Connections []NetworkConnectionStatus
	Peers       []NetworkPeerStatus
	ObservedAt  time.Time
}

type NetworkStatusUpdate struct {
	Snapshot NetworkStatusSnapshot
	Change   NetworkConnectionStatus
}

// NetworkPathClassifier is deliberately independent from routing and dialing.
// Additional evidence sources can be added without changing connection logic.
type NetworkPathClassifier interface {
	Classify(NetworkConnectionMetadata) (NetworkClassification, []string)
}

type conservativeNetworkPathClassifier struct{}

func (conservativeNetworkPathClassifier) Classify(metadata NetworkConnectionMetadata) (NetworkClassification, []string) {
	details := append([]string(nil), metadata.Evidence.Details...)
	if metadata.Evidence.InternetRequired && metadata.Evidence.VerifiedLocalDirect {
		return NetworkClassificationUnknown, append(details, "conflicting verified path evidence")
	}
	if metadata.Evidence.InternetRequired {
		return NetworkClassificationInternet, append(details, "path independently verified as Internet-dependent")
	}
	if metadata.Evidence.VerifiedLocalDirect && metadata.Direct && !metadata.Relayed {
		return NetworkClassificationOffline, append(details, "direct local path independently verified")
	}
	return NetworkClassificationUnknown, append(details, "available connection metadata does not prove Internet or offline-local dependency")
}

type networkStatusObserver struct {
	ctx            context.Context
	cancel         context.CancelFunc
	network        network.Network
	classifier     NetworkPathClassifier
	evidenceSource NetworkConnectionEvidenceSource

	mu          sync.Mutex
	connections map[string]NetworkConnectionStatus
	watchers    map[uint64]chan NetworkStatusUpdate
	nextWatcher uint64
	closed      bool
}

func newNetworkStatusObserver(parent context.Context, n network.Network, source NetworkConnectionEvidenceSource) *networkStatusObserver {
	ctx, cancel := context.WithCancel(parent)
	o := &networkStatusObserver{
		ctx:            ctx,
		cancel:         cancel,
		network:        n,
		classifier:     conservativeNetworkPathClassifier{},
		evidenceSource: source,
		connections:    make(map[string]NetworkConnectionStatus),
		watchers:       make(map[uint64]chan NetworkStatusUpdate),
	}
	if n != nil {
		n.Notify(o)
		for _, conn := range n.Conns() {
			o.recordConnection(conn, NetworkConnectionStateConnected)
		}
	}
	return o
}

func (o *networkStatusObserver) Close() {
	o.mu.Lock()
	if o.closed {
		o.mu.Unlock()
		return
	}
	o.closed = true
	o.mu.Unlock()

	if o.network != nil {
		o.network.StopNotify(o)
	}
	o.cancel()

	o.mu.Lock()
	for id, watcher := range o.watchers {
		delete(o.watchers, id)
		close(watcher)
	}
	o.mu.Unlock()
}

func (o *networkStatusObserver) Snapshot() NetworkStatusSnapshot {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.snapshotLocked(time.Now())
}

func (o *networkStatusObserver) Watch(ctx context.Context) <-chan NetworkStatusUpdate {
	out := make(chan NetworkStatusUpdate, 1)
	o.mu.Lock()
	if o.closed {
		close(out)
		o.mu.Unlock()
		return out
	}
	id := o.nextWatcher
	o.nextWatcher++
	o.watchers[id] = out
	out <- NetworkStatusUpdate{Snapshot: o.snapshotLocked(time.Now())}
	o.mu.Unlock()

	go func() {
		select {
		case <-ctx.Done():
		case <-o.ctx.Done():
		}
		o.mu.Lock()
		if watcher, ok := o.watchers[id]; ok {
			delete(o.watchers, id)
			close(watcher)
		}
		o.mu.Unlock()
	}()
	return out
}

func (o *networkStatusObserver) Connected(_ network.Network, conn network.Conn) {
	o.recordConnection(conn, NetworkConnectionStateConnected)
}

func (o *networkStatusObserver) Disconnected(_ network.Network, conn network.Conn) {
	o.recordConnection(conn, NetworkConnectionStateDisconnected)
}

func (o *networkStatusObserver) Listen(network.Network, ma.Multiaddr)      {}
func (o *networkStatusObserver) ListenClose(network.Network, ma.Multiaddr) {}

func (o *networkStatusObserver) recordConnection(conn network.Conn, state NetworkConnectionState) {
	if conn == nil {
		return
	}
	stats := conn.Stat()
	metadata := NetworkConnectionMetadata{
		PeerID:          conn.RemotePeer().String(),
		ConnectionID:    conn.ID(),
		State:           state,
		LocalMultiaddr:  multiaddrString(conn.LocalMultiaddr()),
		RemoteMultiaddr: multiaddrString(conn.RemoteMultiaddr()),
		Transport:       conn.ConnState().Transport,
		Direct:          !stats.Limited && !isRelayed(conn.LocalMultiaddr()) && !isRelayed(conn.RemoteMultiaddr()),
		Relayed:         stats.Limited || isRelayed(conn.LocalMultiaddr()) || isRelayed(conn.RemoteMultiaddr()),
		ObservedAt:      time.Now(),
	}
	if isVerifiedProximity(metadata) {
		metadata.Evidence.VerifiedLocalDirect = true
		metadata.Evidence.Details = append(metadata.Evidence.Details, "Wesh proximity transport is device-local")
	}
	if o.evidenceSource != nil {
		external := o.evidenceSource.EvidenceForConnection(conn)
		metadata.Evidence.InternetRequired = metadata.Evidence.InternetRequired || external.InternetRequired
		metadata.Evidence.VerifiedLocalDirect = metadata.Evidence.VerifiedLocalDirect || external.VerifiedLocalDirect
		metadata.Evidence.Details = append(metadata.Evidence.Details, external.Details...)
	}
	o.recordMetadata(metadata)
}

func (o *networkStatusObserver) recordMetadata(metadata NetworkConnectionMetadata) {
	classification, details := o.classifier.Classify(metadata)
	status := NetworkConnectionStatus{
		NetworkConnectionMetadata: metadata,
		Classification:            classification,
		EvidenceDetails:           details,
	}

	o.mu.Lock()
	if o.closed {
		o.mu.Unlock()
		return
	}
	if metadata.State == NetworkConnectionStateDisconnected {
		delete(o.connections, metadata.ConnectionID)
	} else {
		o.connections[metadata.ConnectionID] = status
	}
	update := NetworkStatusUpdate{
		Snapshot: o.snapshotLocked(metadata.ObservedAt),
		Change:   status,
	}
	for _, watcher := range o.watchers {
		select {
		case watcher <- update:
		default:
			select {
			case <-watcher:
			default:
			}
			watcher <- update
		}
	}
	o.mu.Unlock()
}

func (o *networkStatusObserver) snapshotLocked(observedAt time.Time) NetworkStatusSnapshot {
	snapshot := NetworkStatusSnapshot{ObservedAt: observedAt}
	for _, connection := range o.connections {
		connection.Evidence.Details = append([]string(nil), connection.Evidence.Details...)
		connection.EvidenceDetails = append([]string(nil), connection.EvidenceDetails...)
		snapshot.Connections = append(snapshot.Connections, connection)
	}
	sort.Slice(snapshot.Connections, func(i, j int) bool {
		if snapshot.Connections[i].PeerID == snapshot.Connections[j].PeerID {
			return snapshot.Connections[i].ConnectionID < snapshot.Connections[j].ConnectionID
		}
		return snapshot.Connections[i].PeerID < snapshot.Connections[j].PeerID
	})
	snapshot.Peers = aggregateNetworkPeers(snapshot.Connections)
	return snapshot
}

func aggregateNetworkPeers(connections []NetworkConnectionStatus) []NetworkPeerStatus {
	byPeer := make(map[string][]NetworkConnectionStatus)
	for _, connection := range connections {
		byPeer[connection.PeerID] = append(byPeer[connection.PeerID], connection)
	}
	peers := make([]NetworkPeerStatus, 0, len(byPeer))
	for peerID, peerConnections := range byPeer {
		peerStatus := NetworkPeerStatus{PeerID: peerID}
		first := peerConnections[0].Classification
		peerStatus.Classification = first
		for _, connection := range peerConnections {
			peerStatus.ConnectionIDs = append(peerStatus.ConnectionIDs, connection.ConnectionID)
			peerStatus.EvidenceDetails = append(peerStatus.EvidenceDetails, connection.EvidenceDetails...)
			if connection.Classification != first {
				peerStatus.Mixed = true
				peerStatus.Classification = NetworkClassificationUnknown
			}
		}
		sort.Strings(peerStatus.ConnectionIDs)
		peerStatus.EvidenceDetails = uniqueStrings(peerStatus.EvidenceDetails)
		peers = append(peers, peerStatus)
	}
	sort.Slice(peers, func(i, j int) bool { return peers[i].PeerID < peers[j].PeerID })
	return peers
}

func isVerifiedProximity(metadata NetworkConnectionMetadata) bool {
	if metadata.Relayed || !metadata.Direct {
		return false
	}
	parts := strings.Split(strings.TrimPrefix(metadata.RemoteMultiaddr, "/"), "/")
	if len(parts) == 0 {
		return false
	}
	switch parts[0] {
	case "ble", "nearby", "mc":
		return metadata.Transport == parts[0]
	default:
		return false
	}
}

func isRelayed(addr ma.Multiaddr) bool {
	return addr != nil && strings.Contains(addr.String(), "/p2p-circuit")
}

func multiaddrString(addr ma.Multiaddr) string {
	if addr == nil {
		return ""
	}
	return addr.String()
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
