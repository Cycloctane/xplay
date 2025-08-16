VERSION=$(shell git describe --tags --always)
LD_FLAGS=-s -w -X main.version=$(VERSION)
OUTPUT_DIR=build

.PHONY: all test windows_x64 linux_x64 linux_static_x64 linux_static_arm64

default: all

test:
	go test --cover -v ./pkg/xspf/ ./internal/router/

windows_x64:
	GOOS=windows GOARCH=amd64 go build -v -o $(OUTPUT_DIR)/xplay_$(VERSION)_$@.exe -ldflags="$(LD_FLAGS)" -trimpath ./cmd/xplay

linux_x64:
	GOOS=linux GOARCH=amd64 go build -v -o $(OUTPUT_DIR)/xplay_$(VERSION)_$@ -ldflags="$(LD_FLAGS)" -trimpath ./cmd/xplay

linux_static_x64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o $(OUTPUT_DIR)/xplay_$(VERSION)_$@ -ldflags="$(LD_FLAGS)" -trimpath ./cmd/xplay

linux_static_arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -v -o $(OUTPUT_DIR)/xplay_$(VERSION)_$@ -ldflags="$(LD_FLAGS)" -trimpath ./cmd/xplay

all: windows_x64 linux_x64 linux_static_x64 linux_static_arm64

clean:
	rm -f $(OUTPUT_DIR)/*
