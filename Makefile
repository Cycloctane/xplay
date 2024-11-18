VERSION=$(shell git describe --tags --always)

WINDOWS=xplay_$(VERSION)_windows_x64.exe
LINUX=xplay_$(VERSION)_linux_x64
LINUX_STATIC=xplay_$(VERSION)_linux_static_x64

LD_FLAGS=-s -w -X main.version=$(VERSION)

.PHONY: all windows linux staticlinux

default: all

windows:
	GOOS=windows GOARCH=amd64 go build -v -o $(WINDOWS) -ldflags="$(LD_FLAGS)" ./cmd/xplay

linux:
	GOOS=linux GOARCH=amd64 go build -v -o $(LINUX) -ldflags="$(LD_FLAGS)" ./cmd/xplay

staticlinux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o $(LINUX_STATIC) -ldflags="$(LD_FLAGS)" ./cmd/xplay

all: windows linux staticlinux
