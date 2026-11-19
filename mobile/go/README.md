# Wesh `go/`

__Please, read the main [`README.md`](../README.md) file first.__

This folder contains most of the go code needed by the project.

## Download

Main binary:

```sh
git clone https://github.com/core-laboratories/weshnet
cd weshnet/mobile/go
make install

# if things go wrong, please run and share the result of
make doctor
```

## Main components

* [`./cmd/...`](./cmd): binaries / entrypoints
  * `cd go; make install`
* [`./pkg/...`](./pkg): packages especially made to be imported by other projects
  * [`./weshprotocol`](./pkg/weshprotocol): the [Wesh Protocol](https://tone.onl/protocol)
  * ...
* [`./internal/`](./internal): internal packages that explain how the project works
  * These packages cannot be imported directly by other projects.
    Open an issue if one should become public.
* [`./framework`](./framework): bridges used by mobile apps

## Usage

### Daemon

```console
foo@bar:~$ wesh daemon
2020-11-16T20:46:31Z	INFO    bty                 metrics listener	{"handler": "/bty/metrics", "listener": "[::]:8888"}
...
```

Now you can interact with the daemon API.

### Mini

Wesh Mini provides a quick terminal client. This section covers the basic commands.

To open Wesh Mini run:

```console
foo@bar:~$ wesh mini
```

Or in the `weshnet/mobile/go/cmd/wesh` folder enter

```console
foo@bar:~$ go run . mini
```

#### Help

`/help` To get all available commands in wesh mini, and what they do.

#### Contact

`/contact share qr` To get a QR code you can scan with a phone. </br>
`/contact share qr` To get a QR code you can scan with a phone. </br>
`/contact accept all` To accept all incoming contact requests.

#### Navigation

For details, use `/keyboard`

Most important: </br>
`ctrl + n` To go to the chats below. </br>
`ctrl + p` To go to the chats above.

Alternatively you can use the following: </br>
`ctrl + arrow down` To go to the chat below. </br>
`ctrl + arrow up` To go to the chat above.

Closing Wesh mini: </br>
`ctrl + c` To exit Wesh mini, or use `esc`.

### Share Invite

```console
foo@bar:~$ wesh share-invite
https://tone.onl/id#contact/oZBLFX8N5b8YvBsaFgLtJDkHPbvkRUVn1kGiKop2cFyUPoqCChedZizRzDZybrd6AosgxhhaBzd4VhN7FviNfyh9qfBRJWc/name=demo
```

### Info

```console
foo@bar:~$ wesh info
{
  "protocol": {
    "process": {
      "version": "v2.439.1",
      "vcsRef": "96e42dd09",
      "uptimeMs": "28",
      "userCpuTimeMs": "530",
      "systemCpuTimeMs": "188",
      "startedAt": "1669913305",
      "rlimitCur": "1048576",
      "numGoroutine": "391",
      "nofile": "34",
      "numCpu": "4",
      "goVersion": "go1.18.8",
      "operatingSystem": "linux",
      "hostName": "8de9e98ac2df",
      "arch": "arm64",
      "rlimitMax": "1048576",
      "pid": "8332",
      "ppid": "7",
      "priority": "20",
      "workingDir": "/Users/foo/go/src/wesh.tech/wesh/go",
      "systemUsername": "foo"
    },
    "p2p": {
      "connectedPeers": "4"
    },
    "orbitdb": {
      "accountMetadata": {
        "progress": "4",
        "maximum": "4"
      }
    }
  },
  "messenger": {
    "protocolInSameProcess": true,
    "db": {
      "accounts": "1"
    }
  }
}
```

### `--help`

```console
foo@bar:~$ wesh -h
USAGE
  wesh [global flags] <subcommand> [flags] [args...]

SUBCOMMANDS
  daemon          start a full Wesh instance (Wesh Protocol + Wesh Messenger)
  account-daemon  start a full Wesh instance (Wesh Account)
  mini            start a terminal-based mini wesh client (some messaging features not compatible with the app)
  version         print software version
  info            display system info
  groupinit       initialize a new multi-member group
  share-invite    share invite link on your terminal or in the dev channel on Discord
  token-server    token server, a basic token server issuer without auth or logging
  repl-server     replication server
  peers           list peers
  export          export messenger data from the specified wesh node
  remote-logs     stream logs from a remote node
  service-key     helper to generate a key for managed services
  push-server     push relay server
  p2p             helper around libp2p
  relay           relay server
  vc-issuer       start a verified credentials issuer service

FLAGS
  -log.file ...                                   log file path (pattern)
  -log.file-filters debug+:bty*,-*.grpc,error+:*  file zapfilter configuration
  -log.filters info+:bty*,-*.grpc,error+:*        stderr zapfilter configuration
  -log.format color                               stderr logging format. can be: json, console, color, light-console, light-color
  -log.ring-filters info+:bty*,-*.grpc,error+:*   ring zapfilter configuration
  -log.ring-size 10                               ring buffer size in MB
  -log.tyber-auto-attach ...                      tyber host addresses to be automatically attached to

ADVANCED
  -log.filters=':default: CUSTOM'  equivalent to -log.filters='info+:bty*,-*.grpc,error+:* CUSTOM'
                                   -> more info at https://github.com/moul/zapfilter

foo@bar:~$ wesh daemon -h
USAGE
  wesh [global flags] daemon [flags]

FLAGS
  -config ...                                                             config file (optional)
  -log.file <store-dir>/logs                                              log file path (pattern)
  -log.file-filters debug+:bty*,-*.grpc,error+:*                          file zapfilter configuration
  -log.filters info+:bty*,-*.grpc,error+:*                                stderr zapfilter configuration
  -log.format color                                                       stderr logging format. can be: json, console, color, light-console, light-color
  -log.ring-filters info+:bty*,-*.grpc,error+:*                           ring zapfilter configuration
  -log.ring-size 10                                                       ring buffer size in MB
  -log.tyber-auto-attach ...                                              tyber host addresses to be automatically attached to
  -metrics.listener ...                                                   Metrics listener, will enable metrics
  -metrics.pedantic false                                                 Enable Metrics pedantic for debug
  -no-qr false                                                            do not print the QR code in terminal on startup
  -no-system-info false                                                   do not print system info on startup
  -node.default-push-token ...                                            base 64 encoded default platform push token
  -node.disable-group-monitor false                                       disable group monitoring
  -node.display-name foo (cli)                                            display name
  -node.init-timeout 1m0s                                                 maximum time allowed for the initialization
  -node.listeners /ip4/127.0.0.1/tcp/9091/grpc                            gRPC API listeners
  -node.no-notif false                                                    disable desktop notifications
  -node.rdv-rotation 24h0m0s                                              rendezvous rotation base for node
  -node.rebuild-db false                                                  reconstruct messenger DB from OrbitDB logs
  -node.restore-export-path ...                                           inits node from a specified export path
  -node.service-insecure false                                            use insecure connection on services
  -p2p.autorelay true                                                     enable autorelay, force private reachability
  -p2p.ble false                                                          if true Bluetooth Low Energy will be enabled
  -p2p.bootstrap :default:                                                ipfs bootstrap node, `:default:` will set ipfs default bootstrap node
  -p2p.dht client                                                         dht mode, can be: `none`, `client`, `server`, `auto`, `autoserver`
  -p2p.dht-randomwalk true                                                if true dht will have randomwalk enable
  -p2p.disable-ipfs-network false                                         disable as much networking feature as possible, useful during development
  -p2p.high-water 200                                                     ConnManager high watermark
  -p2p.ipfs-api-listeners /ip4/127.0.0.1/tcp/5001                         IPFS API listeners
  -p2p.low-water 150                                                      ConnManager low watermark
  -p2p.max-backoff 1h0m0s                                                 maximum p2p backoff duration
  -p2p.mdns true                                                          if true mdns will be enabled
  -p2p.min-backoff 10s                                                    minimum p2p backoff duration
  -p2p.multipeer-connectivity false                                       if true Multipeer Connectivity will be enabled
  -p2p.nearby false                                                       if true Android Nearby will be enabled
  -p2p.poll-interval 1s                                                   how long the discovery system will waits for more peers
  -p2p.rdvp :default:                                                     list of rendezvous point maddr, `:dev:` will add the default devs servers, `:none:` will disable rdvp
  -p2p.static-relays :default:                                            list of static relay maddrs, `:default:` will use statics relays from the config
  -p2p.swarm-announce ...                                                 IPFS announce addrs
  -p2p.swarm-listeners :default:                                          IPFS swarm listeners
  -p2p.swarm-no-announce ...                                              IPFS exclude announce addrs
  -p2p.tinder-dht-driver true                                             if true dht driver will be enable for tinder
  -p2p.tinder-discover true                                               if true enable tinder discovery
  -p2p.tinder-localdiscovery-driver true                                  if true localdiscovery driver will be enable for tinder
  -p2p.tinder-rdvp-driver true                                            if true rdvp driver will be enable for tinder
  -p2p.webui-listener :3999                                               IPFS WebUI listener
  -passphrase ...                                                         optional sharing-link encryption passphrase
  -preset ...                                                             applies various default values, see ADVANCED section below
  -store.dir /Users/foo/Library/Application Support/wesh-tech/wesh  root datastore directory
  -store.inmem false                                                      disable datastore persistence
  -store.shared-dir ...                                                   shared root datastore directory

ADVANCED
  -log.filters=':default: CUSTOM'        equivalent to -log.filters='info+:bty*,-*.grpc,error+:* CUSTOM'
                                         -> more info at https://github.com/moul/zapfilter
  -preset=performance                    better performance: current development defaults
  -preset=anonymity                      better privacy: -tor.mode=required -p2p.mdns=false -p2p.multipeer-connectivity=false -p2p.ble=false -p2p.nearby=false
  -preset=volatile                       similar to performance but optimize for a quick throwable node: -store.inmem=true -p2p.ipfs-api-listeners="" -p2p.swarm-listeners="" -p2p.webui-listener=""
  -p2p.swarm-listeners=:default:,CUSTOM  equivalent to -p2p.swarm-listeners=/ip4/0.0.0.0/tcp/0,/ip6/::/tcp/0,CUSTOM
  -p2p.rdvp=:default:,CUSTOM             equivalent to -p2p.rdvp=/ip4/52.47.79.109/tcp/4040/p2p/12D3KooWKhU...,CUSTOM
                                         -> full list available at https://github.com/core-laboratories/weshnet/tree/master/config)

foo@bar:~$ wesh account-daemon -h
USAGE
  wesh [global flags] account-daemon [flags]

FLAGS
  -log.file ...                                                           log file path (pattern)
  -log.file-filters debug+:bty*,-*.grpc,error+:*                          file zapfilter configuration
  -log.filters info+:bty*,-*.grpc,error+:*                                stderr zapfilter configuration
  -log.format color                                                       stderr logging format. can be: json, console, color, light-console, light-color
  -log.ring-filters info+:bty*,-*.grpc,error+:*                           ring zapfilter configuration
  -log.ring-size 10                                                       ring buffer size in MB
  -log.tyber-auto-attach ...                                              tyber host addresses to be automatically attached to
  -node.account.listeners /ip4/127.0.0.1/tcp/9092/grpc                    gRPC account API listeners
  -node.listeners /ip4/127.0.0.1/tcp/9091/grpc                            gRPC API listeners
  -store.dir /Users/foo/Library/Application Support/wesh-tech/wesh  root datastore directory
  -store.inmem false                                                      disable datastore persistence
  -store.shared-dir ...                                                   shared root datastore directory

ADVANCED
  -log.filters=':default: CUSTOM'  equivalent to -log.filters='info+:bty*,-*.grpc,error+:* CUSTOM'
                                   -> more info at https://github.com/moul/zapfilter

foo@bar:~$ wesh mini -h
USAGE
  wesh [global flags] mini [flags]

FLAGS
  -config ...                                                             config file (optional)
  -log.file <store-dir>/logs                                              log file path (pattern)
  -log.file-filters debug+:bty*,-*.grpc,error+:*                          file zapfilter configuration
  -log.filters info+:bty*,-*.grpc,error+:*                                stderr zapfilter configuration
  -log.format color                                                       stderr logging format. can be: json, console, color, light-console, light-color
  -log.ring-filters info+:bty*,-*.grpc,error+:*                           ring zapfilter configuration
  -log.ring-size 10                                                       ring buffer size in MB
  -log.tyber-auto-attach ...                                              tyber host addresses to be automatically attached to
  -metrics.listener ...                                                   Metrics listener, will enable metrics
  -metrics.pedantic false                                                 Enable Metrics pedantic for debug
  -mini.group ...                                                         group to join, leave empty to create a new group
  -node.default-push-token ...                                            base 64 encoded default platform push token
  -node.disable-group-monitor false                                       disable group monitoring
  -node.display-name foo (cli)                                        display name
  -node.init-timeout 1m0s                                                 maximum time allowed for the initialization
  -node.listeners ...                                                     gRPC API listeners
  -node.no-notif false                                                    disable desktop notifications
  -node.rdv-rotation 24h0m0s                                              rendezvous rotation base for node
  -node.rebuild-db false                                                  reconstruct messenger DB from OrbitDB logs
  -node.remote-addr ...                                                   remote Wesh gRPC API address
  -node.restore-export-path ...                                           inits node from a specified export path
  -node.service-insecure false                                            use insecure connection on services
  -p2p.autorelay true                                                     enable autorelay, force private reachability
  -p2p.ble false                                                          if true Bluetooth Low Energy will be enabled
  -p2p.bootstrap :default:                                                ipfs bootstrap node, `:default:` will set ipfs default bootstrap node
  -p2p.dht client                                                         dht mode, can be: `none`, `client`, `server`, `auto`, `autoserver`
  -p2p.dht-randomwalk true                                                if true dht will have randomwalk enable
  -p2p.disable-ipfs-network false                                         disable as much networking feature as possible, useful during development
  -p2p.high-water 200                                                     ConnManager high watermark
  -p2p.ipfs-api-listeners /ip4/127.0.0.1/tcp/5001                         IPFS API listeners
  -p2p.low-water 150                                                      ConnManager low watermark
  -p2p.max-backoff 1h0m0s                                                 maximum p2p backoff duration
  -p2p.mdns true                                                          if true mdns will be enabled
  -p2p.min-backoff 10s                                                    minimum p2p backoff duration
  -p2p.multipeer-connectivity false                                       if true Multipeer Connectivity will be enabled
  -p2p.nearby false                                                       if true Android Nearby will be enabled
  -p2p.poll-interval 1s                                                   how long the discovery system will waits for more peers
  -p2p.rdvp :default:                                                     list of rendezvous point maddr, `:dev:` will add the default devs servers, `:none:` will disable rdvp
  -p2p.static-relays :default:                                            list of static relay maddrs, `:default:` will use statics relays from the config
  -p2p.swarm-announce ...                                                 IPFS announce addrs
  -p2p.swarm-listeners :default:                                          IPFS swarm listeners
  -p2p.swarm-no-announce ...                                              IPFS exclude announce addrs
  -p2p.tinder-dht-driver true                                             if true dht driver will be enable for tinder
  -p2p.tinder-discover true                                               if true enable tinder discovery
  -p2p.tinder-localdiscovery-driver true                                  if true localdiscovery driver will be enable for tinder
  -p2p.tinder-rdvp-driver true                                            if true rdvp driver will be enable for tinder
  -p2p.webui-listener :3999                                               IPFS WebUI listener
  -preset ...                                                             applies various default values, see ADVANCED section below
  -store.dir /Users/foo/Library/Application Support/wesh-tech/wesh  root datastore directory
  -store.inmem false                                                      disable datastore persistence
  -store.shared-dir ...                                                   shared root datastore directory

ADVANCED
  -log.filters=':default: CUSTOM'        equivalent to -log.filters='info+:bty*,-*.grpc,error+:* CUSTOM'
                                         -> more info at https://github.com/moul/zapfilter
  -preset=performance                    better performance: current development defaults
  -preset=anonymity                      better privacy: -tor.mode=required -p2p.mdns=false -p2p.multipeer-connectivity=false -p2p.ble=false -p2p.nearby=false
  -preset=volatile                       similar to performance but optimize for a quick throwable node: -store.inmem=true -p2p.ipfs-api-listeners="" -p2p.swarm-listeners="" -p2p.webui-listener=""
  -p2p.swarm-listeners=:default:,CUSTOM  equivalent to -p2p.swarm-listeners=/ip4/0.0.0.0/tcp/0,/ip6/::/tcp/0,CUSTOM
  -p2p.rdvp=:default:,CUSTOM             equivalent to -p2p.rdvp=/ip4/52.47.79.109/tcp/4040/p2p/12D3KooWKhU...,CUSTOM
                                         -> full list available at https://github.com/core-laboratories/weshnet/tree/master/config)

foo@bar:~$ wesh version -h
USAGE
  wesh version

ADVANCED
  -log.filters=':default: CUSTOM'  equivalent to -log.filters='info+:bty*,-*.grpc,error+:* CUSTOM'
                                   -> more info at https://github.com/moul/zapfilter

foo@bar:~$ wesh info -h
USAGE
  wesh [global flags] info [flags]

FLAGS
  -config ...                                                             config file (optional)
  -info.anonymize false                                                   anonymize output for sharing
  -info.refresh 0s                                                        refresh every DURATION (0: no refresh)
  -log.file <store-dir>/logs                                              log file path (pattern)
  -log.file-filters debug+:bty*,-*.grpc,error+:*                          file zapfilter configuration
  -log.filters info+:bty*,-*.grpc,error+:*                                stderr zapfilter configuration
  -log.format color                                                       stderr logging format. can be: json, console, color, light-console, light-color
  -log.ring-filters info+:bty*,-*.grpc,error+:*                           ring zapfilter configuration
  -log.ring-size 10                                                       ring buffer size in MB
  -log.tyber-auto-attach ...                                              tyber host addresses to be automatically attached to
  -node.default-push-token ...                                            base 64 encoded default platform push token
  -node.disable-group-monitor false                                       disable group monitoring
  -node.display-name foo (cli)                                        display name
  -node.no-notif false                                                    disable desktop notifications
  -node.rdv-rotation 24h0m0s                                              rendezvous rotation base for node
  -node.rebuild-db false                                                  reconstruct messenger DB from OrbitDB logs
  -node.remote-addr ...                                                   remote Wesh gRPC API address
  -node.restore-export-path ...                                           inits node from a specified export path
  -node.service-insecure false                                            use insecure connection on services
  -p2p.autorelay true                                                     enable autorelay, force private reachability
  -p2p.ble false                                                          if true Bluetooth Low Energy will be enabled
  -p2p.bootstrap :default:                                                ipfs bootstrap node, `:default:` will set ipfs default bootstrap node
  -p2p.dht client                                                         dht mode, can be: `none`, `client`, `server`, `auto`, `autoserver`
  -p2p.dht-randomwalk true                                                if true dht will have randomwalk enable
  -p2p.disable-ipfs-network false                                         disable as much networking feature as possible, useful during development
  -p2p.high-water 200                                                     ConnManager high watermark
  -p2p.ipfs-api-listeners /ip4/127.0.0.1/tcp/5001                         IPFS API listeners
  -p2p.low-water 150                                                      ConnManager low watermark
  -p2p.max-backoff 1h0m0s                                                 maximum p2p backoff duration
  -p2p.mdns true                                                          if true mdns will be enabled
  -p2p.min-backoff 10s                                                    minimum p2p backoff duration
  -p2p.multipeer-connectivity false                                       if true Multipeer Connectivity will be enabled
  -p2p.nearby false                                                       if true Android Nearby will be enabled
  -p2p.poll-interval 1s                                                   how long the discovery system will waits for more peers
  -p2p.rdvp :default:                                                     list of rendezvous point maddr, `:dev:` will add the default devs servers, `:none:` will disable rdvp
  -p2p.static-relays :default:                                            list of static relay maddrs, `:default:` will use statics relays from the config
  -p2p.swarm-announce ...                                                 IPFS announce addrs
  -p2p.swarm-listeners :default:                                          IPFS swarm listeners
  -p2p.swarm-no-announce ...                                              IPFS exclude announce addrs
  -p2p.tinder-dht-driver true                                             if true dht driver will be enable for tinder
  -p2p.tinder-discover true                                               if true enable tinder discovery
  -p2p.tinder-localdiscovery-driver true                                  if true localdiscovery driver will be enable for tinder
  -p2p.tinder-rdvp-driver true                                            if true rdvp driver will be enable for tinder
  -p2p.webui-listener :3999                                               IPFS WebUI listener
  -preset ...                                                             applies various default values, see ADVANCED section below
  -store.dir /Users/foo/Library/Application Support/wesh-tech/wesh  root datastore directory
  -store.inmem false                                                      disable datastore persistence
  -store.shared-dir ...                                                   shared root datastore directory

ADVANCED
  -log.filters=':default: CUSTOM'        equivalent to -log.filters='info+:bty*,-*.grpc,error+:* CUSTOM'
                                         -> more info at https://github.com/moul/zapfilter
  -preset=performance                    better performance: current development defaults
  -preset=anonymity                      better privacy: -tor.mode=required -p2p.mdns=false -p2p.multipeer-connectivity=false -p2p.ble=false -p2p.nearby=false
  -preset=volatile                       similar to performance but optimize for a quick throwable node: -store.inmem=true -p2p.ipfs-api-listeners="" -p2p.swarm-listeners="" -p2p.webui-listener=""
  -p2p.swarm-listeners=:default:,CUSTOM  equivalent to -p2p.swarm-listeners=/ip4/0.0.0.0/tcp/0,/ip6/::/tcp/0,CUSTOM
  -p2p.rdvp=:default:,CUSTOM             equivalent to -p2p.rdvp=/ip4/52.47.79.109/tcp/4040/p2p/12D3KooWKhU...,CUSTOM
                                         -> full list available at https://github.com/core-laboratories/weshnet/tree/master/config)

foo@bar:~$ wesh groupinit -h
USAGE
  wesh groupinit

FLAGS
  -config ...                                     config file (optional)
  -log.file ...                                   log file path (pattern)
  -log.file-filters debug+:bty*,-*.grpc,error+:*  file zapfilter configuration
  -log.filters info+:bty*,-*.grpc,error+:*        stderr zapfilter configuration
  -log.format color                               stderr logging format. can be: json, console, color, light-console, light-color
  -log.ring-filters info+:bty*,-*.grpc,error+:*   ring zapfilter configuration
  -log.ring-size 10                               ring buffer size in MB
  -log.tyber-auto-attach ...                      tyber host addresses to be automatically attached to
  -no-qr false                                    do not print the QR code in terminal
  -passphrase ...                                 optional encryption passphrase

ADVANCED
  -log.filters=':default: CUSTOM'  equivalent to -log.filters='info+:bty*,-*.grpc,error+:* CUSTOM'
                                   -> more info at https://github.com/moul/zapfilter

foo@bar:~$ wesh share-invite -h
USAGE
  wesh [global flags] share-invite [flags]

FLAGS
  -config ...                                                             config file (optional)
  -dev-channel false                                                      post qrcode on dev channel
  -log.file <store-dir>/logs                                              log file path (pattern)
  -log.file-filters debug+:bty*,-*.grpc,error+:*                          file zapfilter configuration
  -log.filters error+:*                                                   stderr zapfilter configuration
  -log.format color                                                       stderr logging format. can be: json, console, color, light-console, light-color
  -log.ring-filters info+:bty*,-*.grpc,error+:*                           ring zapfilter configuration
  -log.ring-size 10                                                       ring buffer size in MB
  -log.tyber-auto-attach ...                                              tyber host addresses to be automatically attached to
  -name ...                                                               override display name
  -no-qr false                                                            do not print the QR code in terminal
  -node.default-push-token ...                                            base 64 encoded default platform push token
  -node.disable-group-monitor false                                       disable group monitoring
  -node.display-name foo (cli)                                        display name
  -node.no-notif false                                                    disable desktop notifications
  -node.rdv-rotation 24h0m0s                                              rendezvous rotation base for node
  -node.rebuild-db false                                                  reconstruct messenger DB from OrbitDB logs
  -node.remote-addr ...                                                   remote Wesh gRPC API address
  -node.restore-export-path ...                                           inits node from a specified export path
  -node.service-insecure false                                            use insecure connection on services
  -p2p.autorelay true                                                     enable autorelay, force private reachability
  -p2p.ble false                                                          if true Bluetooth Low Energy will be enabled
  -p2p.bootstrap :default:                                                ipfs bootstrap node, `:default:` will set ipfs default bootstrap node
  -p2p.dht client                                                         dht mode, can be: `none`, `client`, `server`, `auto`, `autoserver`
  -p2p.dht-randomwalk true                                                if true dht will have randomwalk enable
  -p2p.disable-ipfs-network false                                         disable as much networking feature as possible, useful during development
  -p2p.high-water 200                                                     ConnManager high watermark
  -p2p.ipfs-api-listeners /ip4/127.0.0.1/tcp/5001                         IPFS API listeners
  -p2p.low-water 150                                                      ConnManager low watermark
  -p2p.max-backoff 1h0m0s                                                 maximum p2p backoff duration
  -p2p.mdns true                                                          if true mdns will be enabled
  -p2p.min-backoff 10s                                                    minimum p2p backoff duration
  -p2p.multipeer-connectivity false                                       if true Multipeer Connectivity will be enabled
  -p2p.nearby false                                                       if true Android Nearby will be enabled
  -p2p.poll-interval 1s                                                   how long the discovery system will waits for more peers
  -p2p.rdvp :default:                                                     list of rendezvous point maddr, `:dev:` will add the default devs servers, `:none:` will disable rdvp
  -p2p.static-relays :default:                                            list of static relay maddrs, `:default:` will use statics relays from the config
  -p2p.swarm-announce ...                                                 IPFS announce addrs
  -p2p.swarm-listeners :default:                                          IPFS swarm listeners
  -p2p.swarm-no-announce ...                                              IPFS exclude announce addrs
  -p2p.tinder-dht-driver true                                             if true dht driver will be enable for tinder
  -p2p.tinder-discover true                                               if true enable tinder discovery
  -p2p.tinder-localdiscovery-driver true                                  if true localdiscovery driver will be enable for tinder
  -p2p.tinder-rdvp-driver true                                            if true rdvp driver will be enable for tinder
  -p2p.webui-listener :3999                                               IPFS WebUI listener
  -passphrase ...                                                         optional sharing-link encryption passphrase
  -preset ...                                                             applies various default values, see ADVANCED section below
  -store.dir /Users/foo/Library/Application Support/wesh-tech/wesh  root datastore directory
  -store.inmem false                                                      disable datastore persistence
  -store.shared-dir ...                                                   shared root datastore directory

ADVANCED
  -log.filters=':default: CUSTOM'        equivalent to -log.filters='info+:bty*,-*.grpc,error+:* CUSTOM'
                                         -> more info at https://github.com/moul/zapfilter
  -preset=performance                    better performance: current development defaults
  -preset=anonymity                      better privacy: -tor.mode=required -p2p.mdns=false -p2p.multipeer-connectivity=false -p2p.ble=false -p2p.nearby=false
  -preset=volatile                       similar to performance but optimize for a quick throwable node: -store.inmem=true -p2p.ipfs-api-listeners="" -p2p.swarm-listeners="" -p2p.webui-listener=""
  -p2p.swarm-listeners=:default:,CUSTOM  equivalent to -p2p.swarm-listeners=/ip4/0.0.0.0/tcp/0,/ip6/::/tcp/0,CUSTOM
  -p2p.rdvp=:default:,CUSTOM             equivalent to -p2p.rdvp=/ip4/52.47.79.109/tcp/4040/p2p/12D3KooWKhU...,CUSTOM
                                         -> full list available at https://github.com/core-laboratories/weshnet/tree/master/config)

foo@bar:~$ wesh token-server -h
USAGE
  wesh [global flags] token-server [flags]

FLAGS
  -auth.secret ...                                base64 encoded secret
  -auth.sk ...                                    base64 encoded signature key
  -config ...                                     config file (optional)
  -generate false                                 generate a single token and output it on stdout
  -http.listener 127.0.0.1:8080                   http listener
  -log.file ...                                   log file path (pattern)
  -log.file-filters debug+:bty*,-*.grpc,error+:*  file zapfilter configuration
  -log.filters info+:bty*,-*.grpc,error+:*        stderr zapfilter configuration
  -log.format color                               stderr logging format. can be: json, console, color, light-console, light-color
  -log.ring-filters info+:bty*,-*.grpc,error+:*   ring zapfilter configuration
  -log.ring-size 10                               ring buffer size in MB
  -log.tyber-auto-attach ...                      tyber host addresses to be automatically attached to
  -no-click false                                 disable the login screen and redirect to the next token step directly
  -privacy-policy-url ...                         url of privacy policies
  -svc ...                                        comma separated list of supported services as name@ip:port

ADVANCED
  -log.filters=':default: CUSTOM'  equivalent to -log.filters='info+:bty*,-*.grpc,error+:* CUSTOM'
                                   -> more info at https://github.com/moul/zapfilter

foo@bar:~$ wesh repl-server -h
USAGE
  wesh [global flags] repl-server [flags]

FLAGS
  -config ...                                                             config file (optional)
  -log.file <store-dir>/logs                                              log file path (pattern)
  -log.file-filters debug+:bty*,-*.grpc,error+:*                          file zapfilter configuration
  -log.filters info+:bty*,-*.grpc,error+:*                                stderr zapfilter configuration
  -log.format color                                                       stderr logging format. can be: json, console, color, light-console, light-color
  -log.ring-filters info+:bty*,-*.grpc,error+:*                           ring zapfilter configuration
  -log.ring-size 10                                                       ring buffer size in MB
  -log.tyber-auto-attach ...                                              tyber host addresses to be automatically attached to
  -node.auth-pk ...                                                       Protocol API Authentication Public Key (base64 encoded)
  -node.auth-secret ...                                                   Protocol API Authentication Secret (base64 encoded)
  -node.default-push-token ...                                            base 64 encoded default platform push token
  -node.listeners /ip4/127.0.0.1/tcp/9091/grpc                            gRPC API listeners
  -node.rdv-rotation 24h0m0s                                              rendezvous rotation base for node
  -node.service-insecure false                                            use insecure connection on services
  -p2p.autorelay true                                                     enable autorelay, force private reachability
  -p2p.ble false                                                          if true Bluetooth Low Energy will be enabled
  -p2p.bootstrap :default:                                                ipfs bootstrap node, `:default:` will set ipfs default bootstrap node
  -p2p.dht client                                                         dht mode, can be: `none`, `client`, `server`, `auto`, `autoserver`
  -p2p.dht-randomwalk true                                                if true dht will have randomwalk enable
  -p2p.disable-ipfs-network false                                         disable as much networking feature as possible, useful during development
  -p2p.high-water 200                                                     ConnManager high watermark
  -p2p.ipfs-api-listeners /ip4/127.0.0.1/tcp/5001                         IPFS API listeners
  -p2p.low-water 150                                                      ConnManager low watermark
  -p2p.max-backoff 1h0m0s                                                 maximum p2p backoff duration
  -p2p.mdns true                                                          if true mdns will be enabled
  -p2p.min-backoff 10s                                                    minimum p2p backoff duration
  -p2p.multipeer-connectivity false                                       if true Multipeer Connectivity will be enabled
  -p2p.nearby false                                                       if true Android Nearby will be enabled
  -p2p.poll-interval 1s                                                   how long the discovery system will waits for more peers
  -p2p.rdvp :default:                                                     list of rendezvous point maddr, `:dev:` will add the default devs servers, `:none:` will disable rdvp
  -p2p.static-relays :default:                                            list of static relay maddrs, `:default:` will use statics relays from the config
  -p2p.swarm-announce ...                                                 IPFS announce addrs
  -p2p.swarm-listeners :default:                                          IPFS swarm listeners
  -p2p.swarm-no-announce ...                                              IPFS exclude announce addrs
  -p2p.tinder-dht-driver true                                             if true dht driver will be enable for tinder
  -p2p.tinder-discover true                                               if true enable tinder discovery
  -p2p.tinder-localdiscovery-driver true                                  if true localdiscovery driver will be enable for tinder
  -p2p.tinder-rdvp-driver true                                            if true rdvp driver will be enable for tinder
  -p2p.webui-listener :3999                                               IPFS WebUI listener
  -preset ...                                                             applies various default values, see ADVANCED section below
  -store.dir /Users/foo/Library/Application Support/wesh-tech/wesh  root datastore directory
  -store.inmem false                                                      disable datastore persistence
  -store.shared-dir ...                                                   shared root datastore directory

ADVANCED
  -log.filters=':default: CUSTOM'        equivalent to -log.filters='info+:bty*,-*.grpc,error+:* CUSTOM'
                                         -> more info at https://github.com/moul/zapfilter
  -preset=performance                    better performance: current development defaults
  -preset=anonymity                      better privacy: -tor.mode=required -p2p.mdns=false -p2p.multipeer-connectivity=false -p2p.ble=false -p2p.nearby=false
  -preset=volatile                       similar to performance but optimize for a quick throwable node: -store.inmem=true -p2p.ipfs-api-listeners="" -p2p.swarm-listeners="" -p2p.webui-listener=""
  -p2p.swarm-listeners=:default:,CUSTOM  equivalent to -p2p.swarm-listeners=/ip4/0.0.0.0/tcp/0,/ip6/::/tcp/0,CUSTOM
  -p2p.rdvp=:default:,CUSTOM             equivalent to -p2p.rdvp=/ip4/52.47.79.109/tcp/4040/p2p/12D3KooWKhU...,CUSTOM
                                         -> full list available at https://github.com/core-laboratories/weshnet/tree/master/config)

foo@bar:~$ wesh peers -h
USAGE
  wesh [global flags] peers [flags]

FLAGS
  -config ...                                                             config file (optional)
  -log.file ...                                                           log file path (pattern)
  -log.file-filters debug+:bty*,-*.grpc,error+:*                          file zapfilter configuration
  -log.filters info+:bty*,-*.grpc,error+:*                                stderr zapfilter configuration
  -log.format color                                                       stderr logging format. can be: json, console, color, light-console, light-color
  -log.ring-filters info+:bty*,-*.grpc,error+:*                           ring zapfilter configuration
  -log.ring-size 10                                                       ring buffer size in MB
  -log.tyber-auto-attach ...                                              tyber host addresses to be automatically attached to
  -node.default-push-token ...                                            base 64 encoded default platform push token
  -node.rdv-rotation 24h0m0s                                              rendezvous rotation base for node
  -node.remote-addr ...                                                   remote Wesh gRPC API address
  -node.service-insecure false                                            use insecure connection on services
  -p2p.autorelay true                                                     enable autorelay, force private reachability
  -p2p.ble false                                                          if true Bluetooth Low Energy will be enabled
  -p2p.bootstrap :default:                                                ipfs bootstrap node, `:default:` will set ipfs default bootstrap node
  -p2p.dht client                                                         dht mode, can be: `none`, `client`, `server`, `auto`, `autoserver`
  -p2p.dht-randomwalk true                                                if true dht will have randomwalk enable
  -p2p.disable-ipfs-network false                                         disable as much networking feature as possible, useful during development
  -p2p.high-water 200                                                     ConnManager high watermark
  -p2p.ipfs-api-listeners /ip4/127.0.0.1/tcp/5001                         IPFS API listeners
  -p2p.low-water 150                                                      ConnManager low watermark
  -p2p.max-backoff 1h0m0s                                                 maximum p2p backoff duration
  -p2p.mdns true                                                          if true mdns will be enabled
  -p2p.min-backoff 10s                                                    minimum p2p backoff duration
  -p2p.multipeer-connectivity false                                       if true Multipeer Connectivity will be enabled
  -p2p.nearby false                                                       if true Android Nearby will be enabled
  -p2p.poll-interval 1s                                                   how long the discovery system will waits for more peers
  -p2p.rdvp :default:                                                     list of rendezvous point maddr, `:dev:` will add the default devs servers, `:none:` will disable rdvp
  -p2p.static-relays :default:                                            list of static relay maddrs, `:default:` will use statics relays from the config
  -p2p.swarm-announce ...                                                 IPFS announce addrs
  -p2p.swarm-listeners :default:                                          IPFS swarm listeners
  -p2p.swarm-no-announce ...                                              IPFS exclude announce addrs
  -p2p.tinder-dht-driver true                                             if true dht driver will be enable for tinder
  -p2p.tinder-discover true                                               if true enable tinder discovery
  -p2p.tinder-localdiscovery-driver true                                  if true localdiscovery driver will be enable for tinder
  -p2p.tinder-rdvp-driver true                                            if true rdvp driver will be enable for tinder
  -p2p.webui-listener :3999                                               IPFS WebUI listener
  -peers.refresh 1s                                                       refresh every DURATION (0: no refresh)
  -preset ...                                                             applies various default values, see ADVANCED section below
  -store.dir /Users/foo/Library/Application Support/wesh-tech/wesh  root datastore directory
  -store.inmem false                                                      disable datastore persistence
  -store.shared-dir ...                                                   shared root datastore directory

ADVANCED
  -log.filters=':default: CUSTOM'        equivalent to -log.filters='info+:bty*,-*.grpc,error+:* CUSTOM'
                                         -> more info at https://github.com/moul/zapfilter
  -preset=performance                    better performance: current development defaults
  -preset=anonymity                      better privacy: -tor.mode=required -p2p.mdns=false -p2p.multipeer-connectivity=false -p2p.ble=false -p2p.nearby=false
  -preset=volatile                       similar to performance but optimize for a quick throwable node: -store.inmem=true -p2p.ipfs-api-listeners="" -p2p.swarm-listeners="" -p2p.webui-listener=""
  -p2p.swarm-listeners=:default:,CUSTOM  equivalent to -p2p.swarm-listeners=/ip4/0.0.0.0/tcp/0,/ip6/::/tcp/0,CUSTOM
  -p2p.rdvp=:default:,CUSTOM             equivalent to -p2p.rdvp=/ip4/52.47.79.109/tcp/4040/p2p/12D3KooWKhU...,CUSTOM
                                         -> full list available at https://github.com/core-laboratories/weshnet/tree/master/config)

foo@bar:~$ wesh export -h
USAGE
  wesh [global flags] export [flags]

FLAGS
  -config ...                                                             config file (optional)
  -export-path ...                                                        path of the export tarball
  -log.file ...                                                           log file path (pattern)
  -log.file-filters debug+:bty*,-*.grpc,error+:*                          file zapfilter configuration
  -log.filters info+:bty*,-*.grpc,error+:*                                stderr zapfilter configuration
  -log.format color                                                       stderr logging format. can be: json, console, color, light-console, light-color
  -log.ring-filters info+:bty*,-*.grpc,error+:*                           ring zapfilter configuration
  -log.ring-size 10                                                       ring buffer size in MB
  -log.tyber-auto-attach ...                                              tyber host addresses to be automatically attached to
  -node.default-push-token ...                                            base 64 encoded default platform push token
  -node.disable-group-monitor false                                       disable group monitoring
  -node.display-name foo (cli)                                            display name
  -node.no-notif false                                                    disable desktop notifications
  -node.rdv-rotation 24h0m0s                                              rendezvous rotation base for node
  -node.rebuild-db false                                                  reconstruct messenger DB from OrbitDB logs
  -node.remote-addr ...                                                   remote Wesh gRPC API address
  -node.restore-export-path ...                                           inits node from a specified export path
  -node.service-insecure false                                            use insecure connection on services
  -p2p.autorelay true                                                     enable autorelay, force private reachability
  -p2p.ble false                                                          if true Bluetooth Low Energy will be enabled
  -p2p.bootstrap :default:                                                ipfs bootstrap node, `:default:` will set ipfs default bootstrap node
  -p2p.dht client                                                         dht mode, can be: `none`, `client`, `server`, `auto`, `autoserver`
  -p2p.dht-randomwalk true                                                if true dht will have randomwalk enable
  -p2p.disable-ipfs-network false                                         disable as much networking feature as possible, useful during development
  -p2p.high-water 200                                                     ConnManager high watermark
  -p2p.ipfs-api-listeners /ip4/127.0.0.1/tcp/5001                         IPFS API listeners
  -p2p.low-water 150                                                      ConnManager low watermark
  -p2p.max-backoff 1h0m0s                                                 maximum p2p backoff duration
  -p2p.mdns true                                                          if true mdns will be enabled
  -p2p.min-backoff 10s                                                    minimum p2p backoff duration
  -p2p.multipeer-connectivity false                                       if true Multipeer Connectivity will be enabled
  -p2p.nearby false                                                       if true Android Nearby will be enabled
  -p2p.poll-interval 1s                                                   how long the discovery system will waits for more peers
  -p2p.rdvp :default:                                                     list of rendezvous point maddr, `:dev:` will add the default devs servers, `:none:` will disable rdvp
  -p2p.static-relays :default:                                            list of static relay maddrs, `:default:` will use statics relays from the config
  -p2p.swarm-announce ...                                                 IPFS announce addrs
  -p2p.swarm-listeners :default:                                          IPFS swarm listeners
  -p2p.swarm-no-announce ...                                              IPFS exclude announce addrs
  -p2p.tinder-dht-driver true                                             if true dht driver will be enable for tinder
  -p2p.tinder-discover true                                               if true enable tinder discovery
  -p2p.tinder-localdiscovery-driver true                                  if true localdiscovery driver will be enable for tinder
  -p2p.tinder-rdvp-driver true                                            if true rdvp driver will be enable for tinder
  -p2p.webui-listener :3999                                               IPFS WebUI listener
  -preset ...                                                             applies various default values, see ADVANCED section below
  -store.dir /Users/foo/Library/Application Support/wesh-tech/wesh  root datastore directory
  -store.inmem false                                                      disable datastore persistence
  -store.shared-dir ...                                                   shared root datastore directory

ADVANCED
  -log.filters=':default: CUSTOM'        equivalent to -log.filters='info+:bty*,-*.grpc,error+:* CUSTOM'
                                         -> more info at https://github.com/moul/zapfilter
  -preset=performance                    better performance: current development defaults
  -preset=anonymity                      better privacy: -tor.mode=required -p2p.mdns=false -p2p.multipeer-connectivity=false -p2p.ble=false -p2p.nearby=false
  -preset=volatile                       similar to performance but optimize for a quick throwable node: -store.inmem=true -p2p.ipfs-api-listeners="" -p2p.swarm-listeners="" -p2p.webui-listener=""
  -p2p.swarm-listeners=:default:,CUSTOM  equivalent to -p2p.swarm-listeners=/ip4/0.0.0.0/tcp/0,/ip6/::/tcp/0,CUSTOM
  -p2p.rdvp=:default:,CUSTOM             equivalent to -p2p.rdvp=/ip4/52.47.79.109/tcp/4040/p2p/12D3KooWKhU...,CUSTOM
                                         -> full list available at https://github.com/core-laboratories/weshnet/tree/master/config)

foo@bar:~$ wesh remote-logs -h
USAGE
  wesh [global flags] remote-logs

FLAGS
  -node.remote-addr ...  remote Wesh gRPC API address

ADVANCED
  -log.filters=':default: CUSTOM'  equivalent to -log.filters='info+:bty*,-*.grpc,error+:* CUSTOM'
                                   -> more info at https://github.com/moul/zapfilter

foo@bar:~$ wesh service-key -h
USAGE
  wesh service-key [flags] [command]

FLAGS
  -auth.sk ...  base64 encoded ed25519 private key to convert to a public key

ADVANCED
  -log.filters=':default: CUSTOM'  equivalent to -log.filters='info+:bty*,-*.grpc,error+:* CUSTOM'
                                   -> more info at https://github.com/moul/zapfilter

foo@bar:~$ wesh push-server -h
USAGE
  wesh [global flags] push-server [flags]

FLAGS
  -apns ...                                       Apple's apns certs path, comma-separated
  -fcm ...                                        Firebase's FCM API keys, formatted like app_id:api_key and comma-separated
  -log.file ...                                   log file path (pattern)
  -log.file-filters debug+:bty*,-*.grpc,error+:*  file zapfilter configuration
  -log.filters info+:bty*,-*.grpc,error+:*        stderr zapfilter configuration
  -log.format color                               stderr logging format. can be: json, console, color, light-console, light-color
  -log.ring-filters info+:bty*,-*.grpc,error+:*   ring zapfilter configuration
  -log.ring-size 10                               ring buffer size in MB
  -log.tyber-auto-attach ...                      tyber host addresses to be automatically attached to
  -node.auth-pk ...                               Protocol API Authentication Public Key (base64 encoded)
  -node.auth-secret ...                           Protocol API Authentication Secret (base64 encoded)
  -node.listeners /ip4/127.0.0.1/tcp/9091/grpc    gRPC API listeners
  -push-private-key ...                           Push server private key, base64 formatted

ADVANCED
  -log.filters=':default: CUSTOM'  equivalent to -log.filters='info+:bty*,-*.grpc,error+:* CUSTOM'
                                   -> more info at https://github.com/moul/zapfilter

foo@bar:~$ wesh p2p -h
USAGE
  wesh p2p [command]

SUBCOMMANDS
  generate-id  helper to generate a key identity for libp2p service
  load-id      helper to display peerid or related information from libp2p key from path or stdin

ADVANCED
  -log.filters=':default: CUSTOM'  equivalent to -log.filters='info+:bty*,-*.grpc,error+:* CUSTOM'
                                   -> more info at https://github.com/moul/zapfilter

foo@bar:~$ wesh relay -h
USAGE
  wesh [global flags] relay [flags]

FLAGS
  -key-autogenerate true                                                                                    auto-generate key if none is found
  -key-b64 ...                                                                                              full private key, base 64 encoded
  -key-file peer.key                                                                                        path of the peer key, if none exist, one will be automatically created
  -key-size 2048                                                                                            for RSA key only, specify the bit size of the key
  -key-type Ed25519                                                                                         for RSA key only, specify the bit size of the key
  -log.file ...                                                                                             log file path (pattern)
  -log.file-filters debug+:bty*,-*.grpc,error+:*                                                            file zapfilter configuration
  -log.filters info+:bty*,-*.grpc,error+:*                                                                  stderr zapfilter configuration
  -log.format color                                                                                         stderr logging format. can be: json, console, color, light-console, light-color
  -log.ring-filters info+:bty*,-*.grpc,error+:*                                                             ring zapfilter configuration
  -log.ring-size 10                                                                                         ring buffer size in MB
  -log.tyber-auto-attach ...                                                                                tyber host addresses to be automatically attached to
  -prom-listener :9092                                                                                      listener to expose Prometheus host metric
  -relay-buffsize 2048                                                                                      the size of the relayed connection buffers
  -relay-limit-data 131072                                                                                  Data is the limit of data relayed (on each direction) before resetting the connection. Defaults to 128KB
  -relay-limit-duration 2m0s                                                                                the time limit before resetting a relayed connection
  -relay-max-circuits 2048                                                                                  the maximum number of open relay connections for each peer
  -relay-max-reservations 128                                                                               the maximum number of active relay slots
  -relay-max-reservations-per-asn 32                                                                        the maximum number of reservations origination from the same ASN
  -relay-max-reservations-per-ip 8                                                                          the maximum number of reservations originating from the same IP
  -relay-max-reservations-per-peer 4                                                                        the maximum number of reservations originating from the same peer
  -relay-reservation-ttl 1h0m0s                                                                             ReservationTTL is the duration of a new (or refreshed reservation).
  -swarm-announce ...                                                                                       addrs that will be announce by this server
  -swarm-listeners /ip4/0.0.0.0/tcp/4040,/ip4/0.0.0.0/udp/4141/quic,/ip6/::/tcp/4040,/ip6/::/udp/4040/quic  lists of listeners of (m)addrs separate by a comma

ADVANCED
  -log.filters=':default: CUSTOM'  equivalent to -log.filters='info+:bty*,-*.grpc,error+:* CUSTOM'
                                   -> more info at https://github.com/moul/zapfilter

foo@bar:~$ wesh vc-issuer -h
USAGE
  wesh [global flags] vc-issuer [flags]

FLAGS
  -http.listener ...     http listener
  -http.server-root ...  http server root
  -vc-sk ...             Verifiable Credentials Issuer private key (base64 encoded)

ADVANCED
  -log.filters=':default: CUSTOM'  equivalent to -log.filters='info+:bty*,-*.grpc,error+:* CUSTOM'
                                   -> more info at https://github.com/moul/zapfilter
```

## Other Binaries

### `rdvp`

```console
foo@bar:~$ rdvp -h
USAGE
  rdvp [global flags] <subcommand>

SUBCOMMANDS
  serve
  genkey
  sharekey

foo@bar:~$ rdvp serve -h
USAGE
  rdvp [global flags] serve [flags]

EXAMPLE
  rdvp genkey > rdvp.key
  rdvp serve -pk `cat rdvp.key` -db ./rdvp-store

FLAGS
  -announce ...                                        addrs that will be announce by this server
  -config ...                                          config file (optional)
  -db :memory:                                         rdvp sqlite URN
  -emitter-admin-key ...                               admin key of the emitter-io server
  -emitter-public-addr ...                             if set, will be used to tell the client where to find emitter server
  -emitter-server ...                                  address of the emitter-io server, ie. tcp://127.0.0.1:8080
  -l /ip4/0.0.0.0/tcp/4040,/ip4/0.0.0.0/udp/4141/quic  lists of listeners of (m)addrs separate by a comma
  -log.file stderr                                     if specified, will log everything in JSON into a file and nothing on stderr
  -log.filters info+:*                                 logged namespaces
  -log.format color                                    if specified, will override default log format
  -metrics ...                                         metrics listener, if empty will disable metrics
  -pk ...                                              private key (generated by `rdvp genkey`)

foo@bar:~$ rdvp genkey -h
USAGE
  genkey

FLAGS
  -length 2048          The length (in bits) of the key generated.
  -log.file stderr      if specified, will log everything in JSON into a file and nothing on stderr
  -log.filters info+:*  logged namespaces
  -log.format color     if specified, will override default log format
  -type Ed25519         Type of the private key generated, one of : Ed25519, ECDSA, Secp256k1, RSA

foo@bar:~$ rdvp sharekey -h
USAGE
  rdvp [global flags] sharekey -pk PK

FLAGS
  -log.file stderr      if specified, will log everything in JSON into a file and nothing on stderr
  -log.filters info+:*  logged namespaces
  -log.format color     if specified, will override default log format
  -pk ...               private key (generated by `rdvp genkey`)
```

### `welcomebot`

```console
foo@bar:~$ welcomebot -h
Usage of welcomebot:
  -addr string
	remote 'wesh daemon' address (default "127.0.0.1:9091")
  -display-name string
    	bot's display name (default "foo (Welcome Bot)")
  -log-format string
    	console, json, light-console, light-json, testing (default "console")
  -staff-conversation-link string
    	link of the staff's conversation to join
  -store string
    	store file path (default "./welcomebot.store")
```
