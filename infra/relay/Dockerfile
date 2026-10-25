# This daemon is archived upstream. v0.5.0 is its final release.
FROM golang:1.27.1-alpine3.24 AS builder

ARG GIT_REPOS=https://github.com/libp2p/go-libp2p-relay-daemon.git
ARG GIT_TAG=v0.5.0

RUN apk add --no-cache git \
	&& git clone --depth 1 --branch "${GIT_TAG}" "${GIT_REPOS}" /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -trimpath -o /go/bin/daemon -ldflags="-s -w" ./cmd/libp2p-relay-daemon

FROM alpine:3.24

COPY --from=builder /go/bin/daemon /usr/local/bin/daemon

ENTRYPOINT ["/usr/local/bin/daemon"]
