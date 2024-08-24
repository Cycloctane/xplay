package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/Cycloctane/xplay/internal/mediahandler"
	"github.com/Cycloctane/xplay/internal/router"
)

func validateDir(path string) {
	file, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	if !file.IsDir() {
		panic("not a directory")
	}
}

func main() {
	var output bool
	var listenAddr string
	var listenPort int
	flag.StringVar(&mediahandler.MediaDir, "d", ".", "served directory")
	flag.BoolVar(&mediahandler.Recursive, "r", false, "read directory recursively")
	flag.BoolVar(&output, "w", false, "write xspf to stdout and exit")
	flag.StringVar(&listenAddr, "b", "0.0.0.0", "http server bind address")
	flag.IntVar(&listenPort, "p", 8080, "http server bind port")
	flag.Parse()

	validateDir(mediahandler.MediaDir)
	if output {
		if err := mediahandler.WriteToStdout(); err != nil {
			panic(err)
		}
		return
	}
	addr := net.JoinHostPort(listenAddr, strconv.Itoa(listenPort))
	fmt.Printf("Starting server at http://%s/ ...\n", addr)
	if err := http.ListenAndServe(addr, router.InitRouter()); err != nil {
		panic(err)
	}
}
