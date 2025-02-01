VERSION := $(shell git describe HEAD --exact-match --tags 2> /dev/null || echo "development")
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'oidc-auth-proxy/version.version=$(VERSION)' -X 'oidc-auth-proxy/version.revision=$(REVISION)'

GO ?= go
GOVERSION := $(shell $(GO) version | sed 's|go version go\(.*\)\..* .*|\1|')
GOEXPECTEDVERSION := $(shell grep -E "^go " go.mod | sed 's|go ||')

all: assert_go_version build_linux_amd64

assert_go_version:
	@if \[ "$(GOVERSION)" != "$(GOEXPECTEDVERSION)" \]; then echo "Invalid go version, please use go $(GOEXPECTEDVERSION)"; exit 1; fi

build_linux_amd64:
	@echo "Build for linux amd64..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build -ldflags "$(LDFLAGS)"
