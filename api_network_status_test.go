package weshnet

import (
	"testing"
	"time"

	"github.com/core-laboratories/weshnet/v2/pkg/protocoltypes"
	"github.com/stretchr/testify/require"
)

func TestNetworkStatusSnapshotToProto(t *testing.T) {
	observedAt := time.Unix(123, 456000000)
	snapshot := NetworkStatusSnapshot{
		ObservedAt: observedAt,
		Connections: []NetworkConnectionStatus{{
			NetworkConnectionMetadata: NetworkConnectionMetadata{
				PeerID:          "peer-a",
				ConnectionID:    "connection-a",
				State:           NetworkConnectionStateConnected,
				LocalMultiaddr:  "/ip4/192.0.2.1/tcp/1234",
				RemoteMultiaddr: "/ip4/198.51.100.1/tcp/4321",
				Transport:       "tcp",
				Direct:          true,
				ObservedAt:      observedAt,
			},
			Classification:  NetworkClassificationInternet,
			EvidenceDetails: []string{"test evidence"},
		}},
		Peers: []NetworkPeerStatus{{
			PeerID:          "peer-a",
			Classification:  NetworkClassificationInternet,
			ConnectionIDs:   []string{"connection-a"},
			EvidenceDetails: []string{"test evidence"},
		}},
	}

	status := networkStatusSnapshotToProto(snapshot)
	require.Equal(t, observedAt.UnixMilli(), status.ObservedAtUnixMs)
	require.Len(t, status.Connections, 1)
	require.Equal(t, protocoltypes.NetworkStatusClassification_NetworkStatusClassificationInternet, status.Connections[0].Classification)
	require.Equal(t, protocoltypes.NetworkConnectionState_NetworkConnectionStateConnected, status.Connections[0].State)
	require.Equal(t, "connection-a", status.Connections[0].ConnectionId)
	require.Len(t, status.Peers, 1)
	require.Equal(t, []string{"connection-a"}, status.Peers[0].ConnectionIds)
}
