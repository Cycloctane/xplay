package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/Cycloctane/xplay/internal/mediahandler"
	"github.com/Cycloctane/xplay/internal/router"
)

const (
	defaultPort = 8080
	defaultUser = "xplay"
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
	flag.StringVar(&mediahandler.MediaDir, "d", ".", "served directory")
	flag.BoolVar(&mediahandler.NoTag, "no-tag", false, "do not read media metadata")
	flag.BoolVar(&mediahandler.NoRecursive, "no-recursive", false, "read directory recursively")
	output := flag.Bool("w", false, "write xspf to stdout and exit")
	listenAddr := flag.String("b", "0.0.0.0", "http server bind address")
	listenPort := flag.Int("p", defaultPort, "http server bind port")
	username := flag.String("username", defaultUser, "http basic auth username")
	password := flag.String("password", "", "http basic auth password")
	certFile := flag.String("cert", "", "cert file path for https support")
	keyFile := flag.String("key", "", "cert key path for https support")
	flag.Parse()

	validateDir(mediahandler.MediaDir)
	if *output {
		if err := mediahandler.WriteToStdout(); err != nil {
			panic(err)
		}
		return
	}

	var handler http.Handler
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	if *password != "" {
		handler = router.InitAuthRouter(logger, *username, *password)
	} else {
		handler = router.InitLogRouter(logger)
	}

	addr := net.JoinHostPort(*listenAddr, strconv.Itoa(*listenPort))
	if *certFile != "" && *keyFile != "" {
		logger.Printf("Starting xplay server at https://%s/ ...\n", addr)
		if err := http.ListenAndServeTLS(addr, *certFile, *keyFile, handler); err != nil {
			panic(err)
		}
	} else {
		logger.Printf("Starting xplay server at http://%s/ ...\n", addr)
		if err := http.ListenAndServe(addr, handler); err != nil {
			panic(err)
		}
	}
}
