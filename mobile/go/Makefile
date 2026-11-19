##
## Config
##

GO ?= go
GOPATH ?= $(HOME)/go
GO_TAGS ?= -tags "fts5 sqlite sqlite_unlock_notify"
GO_TEST_OPTS ?= -test.timeout=300s -race -cover -coverprofile=coverage.txt -covermode=atomic $(GO_TAGS)
GO_TEST_PATH ?= ./...
GO_TEST_ENV ?=
CI ?= false

BUILD_DATE ?= `date +%s`
VCS_REF ?= `git rev-parse --short HEAD`
VERSION ?= `cat ../../VERSION`
LDFLAGS ?= -ldflags="-checklinkname=0 -X github.com/core-laboratories/weshnet/mobile/v2/go/pkg/weshversion.VcsRef=$(VCS_REF) -X github.com/core-laboratories/weshnet/mobile/v2/go/pkg/weshversion.Version=$(VERSION)"

# @FIXME(gfanton): on macOS Monterey (12.0.x) we currently need to set the
# environment variable `MallocNanoZone` to 0 to avoid a SIGABRT or SIGSEGV
# see https://github.com/golang/go/issues/49138
MACOS_VERSION=$(shell defaults read /System/Library/CoreServices/SystemVersion.plist 'ProductVersion' 2>/dev/null | sed 's/\.[0-9]$$//')
ifeq ($(MACOS_VERSION),12.0)
GO_TEST_ENV := MallocNanoZone=0 $(GO_TEST_ENV)
endif

ifeq ($(MACOS_VERSION),12.1)
GO_TEST_ENV := MallocNanoZone=0 $(GO_TEST_ENV)
endif

##
## General rules
##

all: help
.PHONY: all


help:
	@echo "Available make commands:"
	@cat Makefile | grep '^[a-z]' | grep -v '=' | cut -d: -f1 | sort | sed 's/^/  /'
.PHONY: help


test: unittest lint tidy
.PHONY: test


unittest: go.unittest
.PHONY: unittest


generate: pb.generate
.PHONY: generate


regenerate: gen.clean generate
.PHONY: regenerate


install: go.install
.PHONY: install


clean:
	rm -rf out/
.PHONY: clean


mini.dev: generate
	go run github.com/githubnemo/CompileDaemon --color=true --build="go install -v" --build-dir="./cmd/wesh" --command="wesh mini"
.PHONY: mini.dev


re: clean generate install
.PHONY: re


tidy: go.tidy
.PHONY: tidy


lint: go.lint
.PHONY: lint


lint.fix: go.fmt
.PHONY: lint.fix

##
## Other rules
##


check-program = $(foreach exec,$(1),$(if $(shell PATH="$(PATH)" which $(exec)),,$(error "No $(exec) in PATH")))

go.tidy: pb.generate
	$(call check-program, $(GO))
	GO111MODULE=on $(GO) mod tidy
.PHONY: go.tidy

go.lint: pb.generate
	$(call check-program, golangci-lint)
	golangci-lint run --timeout=5m $(if $(filter $(CI), false), --verbose) ./...
.PHONY: go.lint

go.unittest: pb.generate
	$(call check-program, $(GO))
	$(GO_TEST_ENV) GO111MODULE=on $(GO) test $(GO_TEST_OPTS) $(GO_TEST_PATH)
.PHONY: go.unittest


go.flappy-tests: pb.generate
	TEST_STABILITY=flappy go run moul.io/testman test -v -test.v -timeout=600s -retry=10 -run ^TestFlappy    ./pkg/weshmessenger
	TEST_STABILITY=flappy go run moul.io/testman test -v -test.v -timeout=600s -retry=10 -run ^TestFlappy    ./pkg/weshbot
	TEST_STABILITY=flappy TEST_SPEED=any GOFLAGS="-tags=fts5,sqlite,sqlite_unlock_notify" go run moul.io/testman test -v -test.v -timeout=600s -retry=10 -run ^TestFlappy ./pkg/weshreplication
	# FIXME: run on other packages too
.PHONY: go.flappy-tests


go.broken-tests: pb.generate
	TEST_STABILITY=broken go run moul.io/testman test -continue-on-error -timeout=1200s -test.timeout=60s -retry=5 -run ^TestBroken ./pkg/weshmessenger
.PHONY: go.broken-tests


go.install: pb.generate
	$(call check-program, $(GO))
	@echo GO111MODULE=on $(GO) install $(GO_TAGS) $(LDFLAGS) -v ./cmd/...
	@GO111MODULE=on $(GO) install $(GO_TAGS) $(LDFLAGS) -v ./cmd/...
.PHONY: go.install


docker.build: pb.generate
	$(call check-program, docker)
	docker build -t weshtech/wesh --platform=linux/amd64 ..
.PHONY: docker.build


docker.fast: pb.generate
	$(call check-program, $(GO) docker)
	@mkdir -p out
	GOOS=linux GOARCH=amd64 GO111MODULE=on $(GO) build $(GO_TAGS) -v -o ./out/wesh-linux-static ./cmd/wesh
	docker run -it --rm -v $(PWD)/out/wesh-linux-static:/bin/wesh $(RUN_OPTS) ubuntu wesh $(ARGS)
.PHONY: docker.fast


print-%:
	@echo $* = $($*)


minimum_go_minor_version = 14
validate-go-version:
	@if [ ! "x`$(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f1`" = "x1" ]; then \
		echo "error: Golang version should be \"1.x\". Please use 1.$(minimum_go_minor_version) or more recent."; \
		exit 1; \
	fi
	@if [ `$(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2` -lt $(minimum_go_minor_version) ]; then \
		echo "error: Golang version is not supported. Please use 1.$(minimum_go_minor_version) or more recent."; \
		exit 1; \
	fi
.PHONY: validate-go-version


##
## Code gen
##


protos_src := $(shell find ../api -type f -name '*.proto')
gen_src := $(protos_src) ../buf.gen.yaml ../buf.gen.tag.yaml Makefile
gen_sum := gen.sum
protoc_opts := -I ../api:`go list -m -mod=mod -f {{.Dir}} github.com/grpc-ecosystem/grpc-gateway`/third_party/googleapis:`go list -m -mod=mod -f {{.Dir}} github.com/gogo/protobuf`:/protobuf
pb.generate: gen.sum validate-go-version
$(gen_sum): $(gen_src)
	$(MAKE) generate_local
.PHONY: pb.generate


generate_local:
	go version
	$(call check-program, shasum buf protoc protoc-gen-go protoc-gen-go-grpc protoc-gen-grpc-gateway protoc-gen-gotag)
	cd ..; buf generate api; cd go
	cd ..; buf generate --template buf.gen.tag.yaml api; cd go
	$(MAKE) go.fmt
	shasum $(gen_src) | sort -k 2 > $(gen_sum).tmp
	mv $(gen_sum).tmp $(gen_sum)
.PHONY: generate_local


go.fmt:
	go run github.com/daixiang0/gci write . \
		--skip-generated -s 'standard,default,prefix(github.com/core-laboratories/weshnet)'
	go run mvdan.cc/gofumpt -w .

.PHONY: go.fmt

pkger.generate:
	$(GO) run github.com/markbates/pkger/cmd/pkger -o go/pkg/assets/
.PHONY: pkger.generate

gen.clean:
	rm -f gen.sum $(wildcard */*/*.pb.go) $(wildcard */*/*pb_test.go) $(wildcard */*/*pb.gw.go)
.PHONY: gen.clean


##
## Doc gen
##


doc.generate: depaware.generate
.PHONY: doc.generate


##
## depaware
##


depaware.generate: ./cmd/*/depaware.txt
.PHONY: depaware.generate

depaware_pkgs = ./cmd/wesh ./cmd/rdvp ./framework/weshbridge ./cmd/welcomebot  ./cmd/wesh-doctor ./cmd/wesh-integration
./cmd/*/depaware.txt: ../go.sum
	@set -xe; for pkg in $(depaware_pkgs); do $(GO) run github.com/tailscale/depaware --update $$pkg; done

depaware.check:
	@set -xe; for pkg in $(depaware_pkgs); do $(GO) run github.com/tailscale/depaware --check $$pkg; done
.PHONY: depaware.check


##
## doctor
##


doctor:
	go run ./cmd/wesh-doctor/
.PHONY: doctor


doctor.verbose:
	go run ./cmd/wesh-doctor/ -v
.PHONY: doctor.verbose


## integration


integration.dev:
	go run ./cmd/wesh-integration
.PHONY: integration.prod


integration.prod:
	go run ./cmd/wesh-integration -integration.testbot=`cat ../config/config.gen.json | jq -r '.wesh.contacts["testbot"].link'`
.PHONY: integration.prod
