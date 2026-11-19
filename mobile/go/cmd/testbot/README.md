# Wesh TestBot

`TestBot` is a simple Wesh bot service.

## Usage (in tests)

Use it in end-to-end tests to check the application, its dependencies, and internet connectivity.

## Usage (as a human)

1. add the first bot as a contact

```txt
█████████████████████████████████████████████████████
█████████████████████████████████████████████████████
████ ▄▄▄▄▄ ██ ▄█▀▄▄█▄▀▄ ▄ ▀▄█▄▄▀ ▀ █▄▄█▀▄█ ▄▄▄▄▄ ████
████ █   █ █▄█  ▄▄█▀ ██    ██▀ █▄▀██▄██ ▀█ █   █ ████
████ █▄▄▄█ ██ ▄▀▄▄▀▄▄▀▄  ▄▄▄  ▄▄ ▄ ▄▀ ▄▄▄█ █▄▄▄█ ████
████▄▄▄▄▄▄▄█ ▀▄▀ ▀ ▀▄▀ ▀ █▄█ █▄█ █ █ █ █▄█▄▄▄▄▄▄▄████
████▄  ▀▄▄▄▄ ▀█▀▄ ██   ▄▄ ▄  █▄█▄▄█▄ ▄█▀█▄▄▀▀█▄█▀████
█████▄█ █▀▄ ▀███▀▄▄█▀  ▄█ █▄ ▄▄  ▀▀▀ ▀█▀ ▄██ ▀█ ▄████
█████ ▄██▀▄█▀▀█▄▄  ▄█▄▄ █ █ ▀▄█ ▄██  ▄▄█▄██▀  █▄▄████
████▄█▄ ▄█▄▀▄██▄██▄██▄▀█▄▄▄▄ █▄▄██      █▄██ █▄ ▄████
████▄   ▄▄▄ █ ▄▀██▄ ▄▄   ▀▀  █▄▄▄▄█▄▄▄▄███▄▄▄ █▄ ████
█████▀▄▀ ▀▄███▀ ▀▄▄ ▄  ▄██▄▀  █▀▄▄ █▄██ █ ▀▄  █ ▄████
█████ █▄ ▄▄▄ ▀▀█▄ ▄▄▄█   ▄▄▄ ▄▄▄▄█▄  █▄▄ ▄▄▄  █ ▄████
█████▄ █ █▄█ █ ▄  ▀▄▄ ▄  █▄█ ██▀ ▀▀▀▀█▄▀ █▄█  █▄▄████
████▀█▀▀▄▄▄ ▄ ▀ ▄ ██ ▀▄▄▄ ▄ ▄▄▄▄▄██ ▄▀▄█    ▄▄██ ████
████▄▄▄▀▀ ▄▀  █ ▄▀ █▀▄▄▀▄▄███▀▀█▄█ ▄ █▀█   ██▀▄ ▄████
█████ █▀▀ ▄ ▄ ▀▄▄  █▄ ▄▄   ▄█▀▄  ▄█ ▄▀▄▀▀████ ▀▀█████
████ █▀█▀ ▄██▄█▄ █▄█▄ ▀▀ █▄██ █▄▄▀██  ██▀▀ ▀█▀█▄ ████
██████▀▄█▄▄ ▄  █ █▄ ▄ █▀ ▄▄ ▀██ ▄██  ▄█▀▀█ ▄█▀█▀ ████
█████▀▀▀ █▄  ███▀█▄ ▄▀▀██▄▄▄█▀█▀██▄▄███▄  ▄ █▄▄  ████
████▄██▄▄█▄█▀  █▄ ▄███ █ ▄▄▄ ██▄▄▄▄ ▄██▀ ▄▄▄ ▄▀█▀████
████ ▄▄▄▄▄ █  ▄ █ ▀ ▀█▄█ █▄█  ▄▄ ▄▀▀▀ ▄█ █▄█  █▄ ████
████ █   █ █▄▄█ ▄ ██ █▄█▄ ▄▄▄▄█▄▄█▄▄ ▄█▀▄▄   ▀█ ▄████
████ █▄▄▄█ █ ▀▄ ▄  █▀█▄ ▀▄█▀▄▀█  ▀▄ ███ ▄██▄ █▄▀█████
████▄▄▄▄▄▄▄█▄▄▄█▄▄▄██▄▄▄▄██▄▄▄▄▄▄██▄▄▄▄██▄█▄██▄▄▄████
█████████████████████████████████████████████████████
▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
https://tone.onl/id#contact/oZBLExQaU8TTubcawgfs8AVcH9FbFV5J1kdEVqgueBaFCMWsADV341pCv26zcHMv3GaTL6UHPzVoHzLoereFQLmBGAh2QKb/name=TestBot
```

1. type `/help`

## Under the hood

* requires two running `wesh daemon` instances.
  * the first one is the "entrypoint"
  * the second one is used for some advanced workflows, i.e., multi-user-group, inviting, etc
* this bot is designed for short testing session, it is often reset, and may forgot your shared conversations

## Suggested Dev Environment

Open 3 terminals:

1. `wesh daemon -node.listeners=/ip4/127.0.0.1/tcp/9092/grpc -store.dir=/tmp/testbot/1`
2. `wesh daemon -node.listeners=/ip4/127.0.0.1/tcp/9093/grpc -store.dir=/tmp/testbot/2`
3. `testbot -debug`

## Deployment

See [../../../tool/deployments/testbot/](../../../tool/deployments/testbot/)
