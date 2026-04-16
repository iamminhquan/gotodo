APP      := gotodo
PKG      := github.com/iamminhquan/gotodo
VERSION  ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS  := -s -w -X $(PKG)/internal/version.Version=$(VERSION)
BUILD    := go build -ldflags=$(LDFLAGS)
OUTPUT   := $(APP)$(if $(filter windows,$(shell go env GOOS)),.exe,)

.PHONY: all build run test lint clean install help

## all: build the binary (default target)
all: build

## build: compile and output the binary
build:
	$(BUILD) -o $(OUTPUT) ./cmd/$(APP)

## run: build and run (pass ARGS="..." to forward arguments)
run: build
	./$(OUTPUT) $(ARGS)

## test: run all unit tests with race detector
test:
	go test -race -count=1 ./...

## lint: vet the source code
lint:
	go vet ./...

## clean: remove build artifacts
clean:
	rm -f $(OUTPUT)

## install: install the binary to GOPATH/bin
install:
	go install -ldflags=$(LDFLAGS) ./cmd/$(APP)

## help: print this help message
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## //'
