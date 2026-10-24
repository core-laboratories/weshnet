# Weshnet Documentation

This directory contains the maintained guides, architecture notes, and generated
protocol reference for Weshnet.

## Guides

- [Installation and development](installation.md)
- [Protocol overview](protocol.md)
- [Network status reporting](network-status.md)
- [Infrastructure deployment](infrastructure.md)
- [Contributing](../CONTRIBUTING.md)
- [Security policy](../SECURITY.md)

## Architecture

- [Architecture overview](architecture/overview.md)
- [Messenger MVP flow](architecture/messenger-mvp.md)
- [Native drivers for Go mobile](architecture/adr-gomobile-ipfs.md)
- [Legacy gRPC bridge](architecture/adr-grpc-bridge.txt)

## API reference

- [Protocol messages and services](apis/protocoltypes.md)
- [OpenAPI description](apis/protocoltypes.swagger.json)
- [Postman collection](apis/weshnet.postman_collection.json)

The files in `docs/apis` are generated from the protocol definitions. Make
changes in [`api/protocol`](../api/protocol) rather than editing generated API
documentation manually.
