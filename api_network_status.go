package weshnet

import (
	"context"

	"github.com/core-laboratories/weshnet/v2/pkg/protocoltypes"
)

func (s *service) GetNetworkStatus(context.Context, *protocoltypes.GetNetworkStatus_Request) (*protocoltypes.GetNetworkStatus_Reply, error) {
	if s.networkStatusObserver == nil {
		return &protocoltypes.GetNetworkStatus_Reply{
			Status: networkStatusSnapshotToProto(NetworkStatusSnapshot{}),
		}, nil
	}
	return &protocoltypes.GetNetworkStatus_Reply{
		Status: networkStatusSnapshotToProto(s.networkStatusObserver.Snapshot()),
	}, nil
}

func (s *service) WatchNetworkStatus(_ *protocoltypes.WatchNetworkStatus_Request, srv protocoltypes.ProtocolService_WatchNetworkStatusServer) error {
	if s.networkStatusObserver == nil {
		return srv.Send(&protocoltypes.WatchNetworkStatus_Reply{
			Status: networkStatusSnapshotToProto(NetworkStatusSnapshot{}),
		})
	}
	updates := s.networkStatusObserver.Watch(srv.Context())
	for {
		select {
		case <-srv.Context().Done():
			return nil
		case update, ok := <-updates:
			if !ok {
				return nil
			}
			reply := &protocoltypes.WatchNetworkStatus_Reply{
				Status: networkStatusSnapshotToProto(update.Snapshot),
			}
			if update.Change.ConnectionID != "" {
				reply.Change = networkConnectionStatusToProto(update.Change)
			}
			if err := srv.Send(reply); err != nil {
				return err
			}
		}
	}
}

func networkStatusSnapshotToProto(snapshot NetworkStatusSnapshot) *protocoltypes.NetworkStatus {
	status := &protocoltypes.NetworkStatus{
		ObservedAtUnixMs: snapshot.ObservedAt.UnixMilli(),
		Connections:      make([]*protocoltypes.NetworkStatus_Connection, 0, len(snapshot.Connections)),
		Peers:            make([]*protocoltypes.NetworkStatus_Peer, 0, len(snapshot.Peers)),
	}
	for _, connection := range snapshot.Connections {
		status.Connections = append(status.Connections, networkConnectionStatusToProto(connection))
	}
	for _, peerStatus := range snapshot.Peers {
		status.Peers = append(status.Peers, &protocoltypes.NetworkStatus_Peer{
			PeerId:         peerStatus.PeerID,
			Classification: networkClassificationToProto(peerStatus.Classification),
			ConnectionIds:  append([]string(nil), peerStatus.ConnectionIDs...),
			Mixed:          peerStatus.Mixed,
			Evidence:       append([]string(nil), peerStatus.EvidenceDetails...),
		})
	}
	return status
}

func networkConnectionStatusToProto(status NetworkConnectionStatus) *protocoltypes.NetworkStatus_Connection {
	return &protocoltypes.NetworkStatus_Connection{
		PeerId:           status.PeerID,
		ConnectionId:     status.ConnectionID,
		State:            networkConnectionStateToProto(status.State),
		LocalMultiaddr:   status.LocalMultiaddr,
		RemoteMultiaddr:  status.RemoteMultiaddr,
		Transport:        status.Transport,
		Direct:           status.Direct,
		Relayed:          status.Relayed,
		Classification:   networkClassificationToProto(status.Classification),
		Evidence:         append([]string(nil), status.EvidenceDetails...),
		ObservedAtUnixMs: status.ObservedAt.UnixMilli(),
	}
}

func networkClassificationToProto(classification NetworkClassification) protocoltypes.NetworkStatusClassification {
	switch classification {
	case NetworkClassificationInternet:
		return protocoltypes.NetworkStatusClassification_NetworkStatusClassificationInternet
	case NetworkClassificationOffline:
		return protocoltypes.NetworkStatusClassification_NetworkStatusClassificationOffline
	default:
		return protocoltypes.NetworkStatusClassification_NetworkStatusClassificationUnknown
	}
}

func networkConnectionStateToProto(state NetworkConnectionState) protocoltypes.NetworkConnectionState {
	switch state {
	case NetworkConnectionStateConnected:
		return protocoltypes.NetworkConnectionState_NetworkConnectionStateConnected
	case NetworkConnectionStateDisconnected:
		return protocoltypes.NetworkConnectionState_NetworkConnectionStateDisconnected
	default:
		return protocoltypes.NetworkConnectionState_NetworkConnectionStateUnknown
	}
}
