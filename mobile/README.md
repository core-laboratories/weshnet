# Weshnet mobile runtime

This directory contains the Account and Messenger services and the native
bridge used by Tone. It is shipped in the Weshnet source release so an app build
does not depend on a separate Wesh source checkout. The bridge lifecycle and
RPC routing live in `go/framework/weshbridge/bridge.go` and `bridge_rpc.go`.

`mobile/` is a separate Go module at
`github.com/core-laboratories/weshnet/mobile/v2`. The generated protocol services use `wesh.*`
package names, with no legacy `wesh.*` service aliases. This is a breaking
change for existing account data and clients. This code incorporates Wesh
work and is not a clean-room service rewrite. See [NOTICE.md](NOTICE.md) and
the Apache-2.0 and MIT license files here. The module depends on the Weshnet
release declared in `go.mod`.
The module's local replacement points to the parent Weshnet source tree. A
checkout of a release tag therefore builds the mobile runtime against the
protocol code from that same tag. Update the `require` version when preparing
the tag; the replacement remains local to this source checkout.

The compatibility source was migrated from the local Wesh fork at
`efb07eb4be3802261901b1e173a8c69aa08e2d5d`, including its working-tree
changes, with Tone's bridge routing and Weshnet RC compatibility updates.

Build and test from this directory:

```sh
GOWORK=off go test ./go/framework/weshbridge
GOWORK=off go build ./go/cmd/tone-wesh
```

When publishing a new Weshnet release, update the Weshnet requirement in this
module and regenerate `go.sum` if the module graph changes. Then update Tone's `tool/wesh/runtime.json`
to the same release version and revision. Build Android, iOS, and desktop
artifacts before cutting over the app.
