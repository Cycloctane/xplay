GO ?= go
OUTPUT_DIR ?= build
VERSION ?= $(shell git describe --tags --always)
BUILD = $(GO) build -v -ldflags="-s -w -X main.version=$(VERSION)" -trimpath

.PHONY: all test clean

default: all

test:
	go test --cover -v ./pkg/xspf/ ./internal/router/

xplay_$(VERSION)_windows_x64.exe:
	GOOS=windows GOARCH=amd64 $(BUILD) -o $(OUTPUT_DIR)/$@ ./cmd/xplay

xplay_$(VERSION)_linux_x64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(BUILD) -o $(OUTPUT_DIR)/$@ ./cmd/xplay

xplay_$(VERSION)_linux_arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(BUILD) -o $(OUTPUT_DIR)/$@ ./cmd/xplay

all: xplay_$(VERSION)_windows_x64.exe xplay_$(VERSION)_linux_x64 xplay_$(VERSION)_linux_arm64

clean:
	rm -f $(OUTPUT_DIR)/*
