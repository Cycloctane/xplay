VERSION=$(shell git describe --tags --always)

WINDOWS=xplay_$(VERSION)_windows_x64.exe
LINUX=xplay_$(VERSION)_linux_x64
LINUX_STATIC=xplay_$(VERSION)_linux_static_x64

LD_FLAGS=-s -w -X main.version=$(VERSION)
OUTPUT_DIR=build

.PHONY: all test windows linux staticlinux

default: all

test:
	go test --cover -v ./pkg/xspf/

windows:
	GOOS=windows GOARCH=amd64 go build -v -o $(OUTPUT_DIR)/$(WINDOWS) -ldflags="$(LD_FLAGS)" ./cmd/xplay

linux:
	GOOS=linux GOARCH=amd64 go build -v -o $(OUTPUT_DIR)/$(LINUX) -ldflags="$(LD_FLAGS)" ./cmd/xplay

staticlinux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o $(OUTPUT_DIR)/$(LINUX_STATIC) -ldflags="$(LD_FLAGS)" ./cmd/xplay

all: windows linux staticlinux

clean:
	rm -f $(OUTPUT_DIR)/*
