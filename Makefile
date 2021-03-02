VERSION:=$(shell git describe --tags --always --abbrev=0 --dirty="-dev")
BUILDFLAGS:=-v -ldflags="-s -w -X main.Version=$(VERSION)"
IMPORT_PATH:=github.com/milgradesec/krypton

.PHONY: all
all: build

.PHONY: build
build:
	go build $(BUILDFLAGS) $(IMPORT_PATH)/cmd/krypton

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: clean
clean:
	go clean
	rm -f krypton.exe