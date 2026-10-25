# Network status reporting

Wesh exposes a read-only view of active libp2p connections through
GetNetworkStatus and WatchNetworkStatus. The observer listens to connection
lifecycle notifications and never dials, closes, upgrades, reprioritizes, or
reroutes a connection.

## Status semantics

Each connection is reported separately, including simultaneous connections to
the same peer. A record includes:

- peer and connection identifiers;
- connected or disconnected state;
- local and remote multiaddresses;
- the transport reported by libp2p;
- whether the libp2p connection is direct or relayed;
- classification evidence; and
- the observation time in Unix milliseconds.

The classifications are intentionally conservative:

- INTERNET means an observer-side integration supplied independently verified
  evidence that the path depends on Internet connectivity.
- OFFLINE means Wesh verified a direct, local-only path. The built-in
  classifier recognizes Wesh proximity transports (ble, nearby, and mc).
- UNKNOWN means the available evidence cannot prove either condition.

A TCP, QUIC, or WebSocket transport does not by itself prove Internet use.
Likewise, public/private IP scope and operating-system connectivity state are
not sufficient evidence. Relay connections remain UNKNOWN unless an
observer-side evidence source independently proves Internet dependency. Wesh
does not infer or report an end-to-end route from the relay address.

Peer summaries retain all connection identifiers. If a peer has connections
with different classifications, its summary is marked mixed and its
classification is UNKNOWN; clients should inspect the individual connection
records.

OFFLINE describes the observed network path. It does not mean that the
recipient is present, that a message can currently be delivered, or that a
message has been read.

## API usage

Fetch the current active-connection snapshot:

    status, err := client.GetNetworkStatus(ctx, &protocoltypes.GetNetworkStatus_Request{})
    if err != nil {
        return err
    }
    for _, connection := range status.GetStatus().GetConnections() {
        fmt.Printf("%s: %s\n", connection.GetPeerId(), connection.GetClassification())
    }

Subscribe to an initial snapshot and later lifecycle changes:

    watch, err := client.WatchNetworkStatus(ctx, &protocoltypes.WatchNetworkStatus_Request{})
    if err != nil {
        return err
    }
    for {
        update, err := watch.Recv()
        if err != nil {
            return err
        }
        // update.Status is the current active snapshot.
        // update.Change contains the connected or disconnected event.
    }

Cancel the context to stop a watcher. The server removes the subscription and
does not leave a polling goroutine behind.

## Messenger display guidance

Display the aggregate label only when it is useful, and retain access to
per-connection details:

- show **Internet** for INTERNET;
- show **Offline path** or **Nearby** for OFFLINE, avoiding wording that
  implies recipient presence;
- show **Unknown network path** for UNKNOWN;
- for mixed, show **Mixed network paths** and list individual connections.

Do not replace delivery receipts, peer presence, or application connectivity
state with this classification.

## Limitations and extension

libp2p does not expose reliable Internet-dependency evidence for ordinary
connections. Those connections therefore remain UNKNOWN. Platforms that can
prove such a fact may provide a NetworkConnectionEvidenceSource when creating
the Wesh service. The source is observational only and must not alter the
connection. No external service, telemetry, Luna mode, or routing behavior is
included in this feature.
