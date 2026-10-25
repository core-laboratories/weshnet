# Wesh Network Toolkit

Weshnet is a messenger-independent toolkit for private, resilient,
peer-to-peer communication. Applications can embed the Wesh protocol and
communicate directly when peers are nearby, through configured rendezvous and
relay infrastructure when they are remote, and across intermittent network
connections.

The runtime is implemented in Go. Its protocol surface is exposed through gRPC,
so it can be integrated with Go, mobile bindings, and any language supported by
the gRPC toolchain.

## Capabilities

- Encrypted one-to-one and multi-member group communication
- Account, contact, device, and group lifecycle management
- Peer discovery through local and remote rendezvous mechanisms
- Transport support for libp2p and platform proximity drivers
- Persistent event logs and message replication
- Configurable relay and rendezvous infrastructure
- Rotating rendezvous seeds for retiring previous discovery material
- Generated protobuf, gRPC, and OpenAPI interfaces

## Requirements

- Go version declared in [`go.mod`](go.mod)
- Git and Make
- Repository tools installed through `asdf` for generation and lint tasks
- Docker when regenerating protocol artifacts or running infrastructure

See the [installation guide](docs/installation.md) for platform-specific setup.

## Install

Use the versioned module path because Weshnet is a Go module at major version 2:

```sh
go get github.com/core-laboratories/weshnet/v2@latest
```

Import the root package:

```go
import "github.com/core-laboratories/weshnet/v2"
```

## Minimal service setup

Applications normally create a persistent service client and configure relay or
rendezvous addresses for their environment:

```go
package main

import (
	"log"

	weshnet "github.com/core-laboratories/weshnet/v2"
)

func main() {
	client, err := weshnet.NewPersistentServiceClient("./wesh-data")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Use client.ClientConn() with generated protocol clients.
}
```

For production use, handle shutdown errors, configure logging, and provide the
network options required by your deployment.

## Repository layout

| Path | Purpose |
| --- | --- |
| `api/go-internal` | Internal protobuf definitions used by the runtime |
| `api/protocol` | Public protocol and service definitions |
| `docs` | Guides, architecture records, and generated API documentation |
| `infra` | Rendezvous and relay deployment resources |
| `internal` | Private implementation packages |
| `pkg` | Reusable protocol, transport, discovery, and utility packages |
| `tool` | Development and benchmarking tools |

## Development

Run the standard validation suite:

```sh
make test
```

Useful targets include:

```sh
make unittest
make lint
make lint.fix
make generate
```

Generation updates checked-in protobuf artifacts. Review generated changes
before committing them.

## Documentation

- [Documentation index](docs/README.md)
- [Installation and development](docs/installation.md)
- [Protocol overview](docs/protocol.md)
- [Network status reporting](docs/network-status.md)
- [Infrastructure deployment](docs/infrastructure.md)
- [Architecture](docs/architecture/overview.md)
- [Protocol API reference](docs/apis/protocoltypes.md)

## Compatibility

The `/v2` suffix is part of the public Go module identity and must remain in
imports. Protocol or storage changes may require coordinated upgrades across
peers. Review release notes and migration guidance before upgrading production
deployments.

## Contributing and security

Read the [contribution guide](CONTRIBUTING.md) before opening a pull request.
Report vulnerabilities privately according to the
[security policy](SECURITY.md).

For regular bug reports and feature requests, use the
[GitHub issue tracker](https://github.com/core-laboratories/weshnet/issues).

## License

Weshnet is available under the [MIT License](LICENSE-MIT).
