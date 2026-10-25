# Infrastructure Deployment

This guide covers the optional rendezvous, Emitter, and circuit-relay services
included under [`../infra`](../infra). These services improve remote peer
discovery and connectivity; they do not decrypt group traffic.

## Components

| Component | Purpose | Default exposed port |
| --- | --- | --- |
| Rendezvous point | Advertises peer discovery records | TCP and UDP `4040` |
| Emitter | Delivers rendezvous presence notifications | TCP `9494` |
| Circuit relay | Relays traffic when direct dialing fails | TCP and UDP `6363` |

## Prerequisites

- A Linux host or another Docker-capable system
- Docker Engine with the Compose plugin
- Make
- A stable public IP address or DNS name
- Firewall and router access for the published TCP and UDP ports

Protect generated keys and service credentials. Do not commit production values
to the repository.

## Rendezvous and Emitter

The Compose project is in [`../infra/rdvp`](../infra/rdvp).

### Configure the service

Copy the environment template and set `RDVP_IMAGE` to an image that provides a
compatible `rdvp` executable. Pin production images to a digest; this repository
does not silently depend on a messenger-branded container image.

```sh
cd infra/rdvp
cp .env.example .env
```

Then generate a rendezvous private key:

```sh
make genkey
```

The command writes `RDVP_PK` to the local `.env` file. Add the public address and
Emitter credentials to that file:

```dotenv
ANNOUNCE_SERVER=203.0.113.10
EMITTER_LICENSE=replace-with-emitter-license
EMITTER_SECRET_KEY=replace-with-emitter-secret
```

Emitter `v3.1` is the latest released server version. Generate an Emitter license
according to its documentation. Restrict permissions on `.env` and back it up
through the deployment's secret-management system. The `.env` file is ignored by
Git; only `.env.example` belongs in source control.

### Start and inspect the service

```sh
make up
make ps
make logs
```

Allow inbound TCP and UDP traffic on port `4040`, and TCP traffic on `9494`.
Confirm that the announced multiaddresses contain the public address rather than
a container or private-network address.

Stop the stack with:

```sh
make down
```

## Circuit relay

The relay Compose project is in [`../infra/relay`](../infra/relay). It uses
`go-libp2p-relay-daemon` `v0.5.0`, the final upstream release. The upstream
project is archived, so treat this stack as a compatibility deployment and plan
to replace it with a maintained application embedding the current libp2p relay
service.

The compatibility image currently builds with Go `1.27.1` on Alpine `3.24` and
uses Alpine `3.24` at runtime. These base versions are deliberately explicit so
upstream image changes cannot silently alter a deployment build.

### Configure the relay

Replace `<ip_address>` in `infra/relay/config.json` with the relay's public
address. Keep TCP and QUIC announce addresses consistent with the published
ports.

Build and start the image:

```sh
cd infra/relay
make build
make up
make ps
make logs
```

Allow inbound TCP and UDP traffic on port `6363`. The `data` directory contains
the persistent relay identity and should be protected and backed up.

## Configure an embedding application

Pass the advertised libp2p multiaddresses when creating the service client:

```go
client, err := weshnet.NewPersistentServiceClient(
	"./wesh-data",
	weshnet.WithP2PRdvpMaddrs([]string{
		"/ip4/203.0.113.10/udp/4040/quic-v1/p2p/12D3KooWExample",
	}),
	weshnet.WithP2PStaticRelays([]string{
		"/ip4/203.0.113.20/udp/6363/quic-v1/p2p/12D3KooWExample",
	}),
)
if err != nil {
	return err
}
defer client.Close()
```

Replace the example peer IDs and addresses with values reported by the deployed
services.

## Production checklist

- Pin container images to immutable digests and regularly review them.
- Store secrets outside source control.
- Restrict administrative and metrics endpoints.
- Monitor container health, restarts, disk usage, and open connections.
- Back up persistent identity data and test restoration.
- Rotate credentials through a documented procedure.
- Keep host and container runtimes patched.
- Test TCP and QUIC access from outside the hosting network.

## Troubleshooting

### Peers receive private addresses

Check `ANNOUNCE_SERVER` and relay `AnnounceAddrs`. They must describe addresses
reachable by clients.

### QUIC does not connect

Confirm UDP forwarding and firewall rules. TCP connectivity does not prove that
UDP is available.

### Relay identity changes after restart

Verify that the relay data volume is persistent and writable by the container.

### Discovery works only locally

Inspect rendezvous and Emitter logs, validate the advertised peer ID, and confirm
that clients use the same deployment addresses.
