# Installation and Development

This guide prepares a local Weshnet checkout for development, testing, and code
generation.

## Prerequisites

- Git
- Make
- Go version declared in [`../go.mod`](../go.mod)
- `asdf` for the repository-pinned development tools
- Docker for protobuf generation and infrastructure services

On macOS, install the command-line developer tools first:

```sh
xcode-select --install
```

Install Homebrew and `asdf` using their official instructions, then install the
repository tools:

```sh
brew install asdf gpg
```

On Ubuntu or Debian, install the base packages before following the official
`asdf` installation guide:

```sh
sudo apt update
sudo apt install build-essential curl git make
```

## Clone the repository

```sh
git clone https://github.com/core-laboratories/weshnet.git
cd weshnet
```

Install the pinned tools after cloning or whenever `.tool-versions` changes:

```sh
make asdf.install_tools
```

Confirm that Go can resolve the module graph:

```sh
go mod download
```

## Validate the checkout

Run the complete local validation target:

```sh
make test
```

The target runs unit tests, lint checks, and module tidying checks. Individual
commands are available when iterating:

```sh
make unittest
make lint
make lint.fix
```

## Generate protocol code

Protocol definitions live under `api`. Regenerate checked-in Go, gRPC, and
related outputs after changing a `.proto` file:

```sh
make generate
```

Generation may use Docker and the versions configured by the repository. Never
edit generated `.pb.go`, gateway, or API reference files by hand. Review the
resulting diff and run `make test` afterward.

## Use Weshnet from another Go module

Install the latest tagged version:

```sh
go get github.com/core-laboratories/weshnet/v2@latest
```

Import packages with the complete major-version suffix:

```go
import weshnet "github.com/core-laboratories/weshnet/v2"
```

Go requires `/v2` for modules at major version 2 or later.

## Troubleshooting

### The Go version is rejected

Install the version declared by `go.mod` or allow the Go toolchain to obtain a
compatible version. Confirm the active toolchain with `go version`.

### Generated files differ after checkout

Verify the pinned tools, Docker version, and repository branch. Regenerate once,
then inspect the diff. Do not discard unrelated local work.

### Module downloads fail

Check proxy and private-module settings with `go env GOPROXY GOPRIVATE`, verify
network access, and retry `go mod download`.

### Native or mobile builds fail

Confirm that the platform SDK, compiler, and Go mobile tooling match the target
platform. The generic Go test suite does not install Android or iOS SDKs.

## Next steps

- Read the [protocol overview](protocol.md).
- Review the [architecture documentation](architecture/overview.md).
- Configure optional [network infrastructure](infrastructure.md).
- Follow the [contribution guide](../CONTRIBUTING.md) before submitting changes.
