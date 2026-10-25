package weshnet

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConservativeNetworkPathClassifier(t *testing.T) {
	classifier := conservativeNetworkPathClassifier{}
	tests := []struct {
		name     string
		metadata NetworkConnectionMetadata
		want     NetworkClassification
	}{
		{
			name: "verified local direct",
			metadata: NetworkConnectionMetadata{
				Direct:   true,
				Evidence: NetworkConnectionEvidence{VerifiedLocalDirect: true},
			},
			want: NetworkClassificationOffline,
		},
		{
			name: "Internet evidenced",
			metadata: NetworkConnectionMetadata{
				Direct:    true,
				Transport: "tcp",
				Evidence:  NetworkConnectionEvidence{InternetRequired: true},
			},
			want: NetworkClassificationInternet,
		},
		{
			name: "relay is not local proof",
			metadata: NetworkConnectionMetadata{
				Relayed:  true,
				Evidence: NetworkConnectionEvidence{VerifiedLocalDirect: true},
			},
			want: NetworkClassificationUnknown,
		},
		{
			name: "transport and public address are ambiguous",
			metadata: NetworkConnectionMetadata{
				Direct:          true,
				Transport:       "quic-v1",
				RemoteMultiaddr: "/ip4/198.51.100.1/udp/443/quic-v1",
			},
			want: NetworkClassificationUnknown,
		},
		{
			name: "conflicting evidence is ambiguous",
			metadata: NetworkConnectionMetadata{
				Direct: true,
				Evidence: NetworkConnectionEvidence{
					InternetRequired:    true,
					VerifiedLocalDirect: true,
				},
			},
			want: NetworkClassificationUnknown,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, evidence := classifier.Classify(test.metadata)
			require.Equal(t, test.want, got)
			require.NotEmpty(t, evidence)
		})
	}
}

func TestVerifiedProximityClassification(t *testing.T) {
	metadata := NetworkConnectionMetadata{
		Direct:          true,
		Transport:       "nearby",
		RemoteMultiaddr: "/nearby/12D3KooWTest",
	}
	require.True(t, isVerifiedProximity(metadata))
	metadata.Relayed = true
	require.False(t, isVerifiedProximity(metadata))
}

func TestNetworkStatusMultipleConnectionsDoNotHideMixedStatus(t *testing.T) {
	observer := newNetworkStatusObserver(context.Background(), nil, nil)
	t.Cleanup(observer.Close)
	now := time.Unix(1700000000, 0)

	observer.recordMetadata(NetworkConnectionMetadata{
		PeerID:       "peer-a",
		ConnectionID: "offline",
		State:        NetworkConnectionStateConnected,
		Direct:       true,
		ObservedAt:   now,
		Evidence:     NetworkConnectionEvidence{VerifiedLocalDirect: true},
	})
	observer.recordMetadata(NetworkConnectionMetadata{
		PeerID:       "peer-a",
		ConnectionID: "unknown",
		State:        NetworkConnectionStateConnected,
		Direct:       true,
		ObservedAt:   now,
	})

	snapshot := observer.Snapshot()
	require.Len(t, snapshot.Connections, 2)
	require.Len(t, snapshot.Peers, 1)
	require.True(t, snapshot.Peers[0].Mixed)
	require.Equal(t, NetworkClassificationUnknown, snapshot.Peers[0].Classification)
	require.Equal(t, []string{"offline", "unknown"}, snapshot.Peers[0].ConnectionIDs)
}

func TestNetworkStatusReconnect(t *testing.T) {
	observer := newNetworkStatusObserver(context.Background(), nil, nil)
	t.Cleanup(observer.Close)
	now := time.Unix(1700000000, 0)
	connected := NetworkConnectionMetadata{
		PeerID:       "peer-a",
		ConnectionID: "connection-1",
		State:        NetworkConnectionStateConnected,
		ObservedAt:   now,
	}
	observer.recordMetadata(connected)
	require.Len(t, observer.Snapshot().Connections, 1)

	disconnected := connected
	disconnected.State = NetworkConnectionStateDisconnected
	disconnected.ObservedAt = now.Add(time.Second)
	observer.recordMetadata(disconnected)
	require.Empty(t, observer.Snapshot().Connections)

	reconnected := connected
	reconnected.ConnectionID = "connection-2"
	reconnected.ObservedAt = now.Add(2 * time.Second)
	observer.recordMetadata(reconnected)
	snapshot := observer.Snapshot()
	require.Len(t, snapshot.Connections, 1)
	require.Equal(t, "connection-2", snapshot.Connections[0].ConnectionID)
}

func TestNetworkStatusConcurrentUpdates(t *testing.T) {
	observer := newNetworkStatusObserver(context.Background(), nil, nil)
	t.Cleanup(observer.Close)

	const connectionCount = 100
	var wg sync.WaitGroup
	for i := 0; i < connectionCount; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			observer.recordMetadata(NetworkConnectionMetadata{
				PeerID:       fmt.Sprintf("peer-%02d", index%5),
				ConnectionID: fmt.Sprintf("connection-%03d", index),
				State:        NetworkConnectionStateConnected,
				ObservedAt:   time.Unix(int64(index), 0),
			})
		}(i)
	}
	wg.Wait()

	snapshot := observer.Snapshot()
	require.Len(t, snapshot.Connections, connectionCount)
	require.Len(t, snapshot.Peers, 5)
}

func TestNetworkStatusWatcherCancellation(t *testing.T) {
	observer := newNetworkStatusObserver(context.Background(), nil, nil)
	t.Cleanup(observer.Close)
	ctx, cancel := context.WithCancel(context.Background())
	updates := observer.Watch(ctx)
	<-updates

	cancel()
	require.Eventually(t, func() bool {
		observer.mu.Lock()
		defer observer.mu.Unlock()
		return len(observer.watchers) == 0
	}, time.Second, 10*time.Millisecond)

	_, ok := <-updates
	require.False(t, ok)
}

func TestNetworkStatusWatcherReceivesDisconnectedChange(t *testing.T) {
	observer := newNetworkStatusObserver(context.Background(), nil, nil)
	t.Cleanup(observer.Close)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	updates := observer.Watch(ctx)
	<-updates

	metadata := NetworkConnectionMetadata{
		PeerID:       "peer-a",
		ConnectionID: "connection-1",
		State:        NetworkConnectionStateConnected,
		ObservedAt:   time.Unix(1700000000, 0),
	}
	observer.recordMetadata(metadata)
	require.Equal(t, NetworkConnectionStateConnected, (<-updates).Change.State)

	metadata.State = NetworkConnectionStateDisconnected
	observer.recordMetadata(metadata)
	update := <-updates
	require.Equal(t, NetworkConnectionStateDisconnected, update.Change.State)
	require.Empty(t, update.Snapshot.Connections)
}
