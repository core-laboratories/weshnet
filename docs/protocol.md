# Wesh Protocol

Wesh is an application-independent protocol for private peer-to-peer messaging.
The repository contains the protocol services, encrypted event stores, discovery
logic, transport integration, and generated interfaces required by an embedding
application.

## Design goals

- Keep application presentation and product policy outside the networking layer.
- Encrypt message and metadata flows using group-scoped cryptographic material.
- Support one-to-one and multi-member conversations.
- Continue operating across intermittent or local-only network conditions.
- Allow discovery and transport implementations to be configured by the host.
- Expose stable language-neutral interfaces through protobuf and gRPC.

## Core concepts

### Account

An account owns its account group and cryptographic identity. The account log
tracks contacts, joined groups, devices, and account-level protocol events.

### Group

A group is the unit of membership, encryption, metadata, and message exchange.
One-to-one conversations use a group shared by two contacts. Multi-member groups
maintain membership events in their metadata store.

### Stores and events

Group metadata and application messages are written to replicated OrbitDB-based
logs. Applications subscribe to protocol streams and interpret the typed events
defined in protobuf.

### Discovery

The Tinder discovery layer combines configured rendezvous mechanisms and
proximity transports. Rendezvous points help peers find each other without
embedding a centralized contact directory in the protocol.

### Rendezvous rotation

A group can register a replacement rendezvous seed. Weshnet derives rotating
discovery topics from the seed and the current time period, keeps the previous
rotation available during a bounded grace period, and removes retired seeds when
requested. Applications should distribute replacement seeds through an
authenticated group channel before retiring the previous seed.

### Replication and relays

Connected peers replicate authorized group logs directly. Relay services help
when direct connectivity is unavailable; they do not replace group-level access
control or encryption.

## API surface

Public definitions are located in [`../api/protocol`](../api/protocol). Generated
Go packages are checked into [`../pkg`](../pkg), and the generated protocol
reference is available in [`apis/protocoltypes.md`](apis/protocoltypes.md).

The main service implements the generated `ProtocolServiceServer`. Embedding
applications can use a persistent client and generated gRPC clients rather than
calling implementation details directly.

## Typical contact flow

1. Alice enables contact requests and obtains a shareable reference.
2. Bob imports the reference and sends a contact request.
3. Alice accepts the request.
4. Both peers receive the membership events for their one-to-one group.
5. The peers exchange and replicate encrypted group messages.

See the [Messenger MVP flow](architecture/messenger-mvp.md) for a detailed
sequence diagram.

## Operational considerations

- Persist service data in a protected application directory.
- Treat account exports, seeds, and private keys as secrets.
- Configure relay and rendezvous addresses appropriate to the deployment.
- Propagate cancellation and close clients cleanly during application shutdown.
- Coordinate protocol upgrades across communicating peers.
- Validate untrusted protobuf and transport input at application boundaries.

## Related documentation

- [Architecture overview](architecture/overview.md)
- [Infrastructure deployment](infrastructure.md)
- [Protocol API reference](apis/protocoltypes.md)
- [Security policy](../SECURITY.md)
