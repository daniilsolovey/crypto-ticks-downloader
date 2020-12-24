NAME = $(notdir $(shell pwd))

VERSION = $(shell printf "%s.%s" \
	$$(git rev-list --count HEAD) \
	$$(git rev-parse --short HEAD) \
)

# could be "..."
TARGET =

GOFLAGS = GO111MODULE=on CGO_ENABLED=0

version:
	@echo $(VERSION)

test:
	$(GOFLAGS) go test -failfast -v ./$(TARGET)

get:
	$(GOFLAGS) go get -v -d

build:
	$(GOFLAGS) go build \
		 -ldflags="-s -w -X main.version=$(VERSION)" \
		 -gcflags="-trimpath=$(GOPATH)" \
		 ./$(TARGET)

all: build
