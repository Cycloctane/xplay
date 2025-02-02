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

var version = "dev"

var (
	showVersion = flag.Bool("version", false, "print version and exit")
	output      = flag.Bool("w", false, "write xspf to stdout and exit")
	listenPort  = flag.Int("p", defaultPort, "http server bind port")
	listenAddr  = flag.String("b", "0.0.0.0", "http server bind address")
	username    = flag.String("username", defaultUser, "http basic auth username")
	password    = flag.String("password", "", "http basic auth password")
	certFile    = flag.String("ssl-cert", "", "cert file path for https support")
	keyFile     = flag.String("ssl-key", "", "cert key path for https support")
)

func init() {
	flag.StringVar(&mediahandler.MediaDir, "d", ".", "served directory")
	flag.BoolVar(&mediahandler.NoTag, "no-tag", false, "do not read media metadata")
	flag.BoolVar(&mediahandler.NoRecursive, "no-recursive", false, "do not read directory recursively")
}

func validateDir(path string, logger *log.Logger) {
	file, err := os.Stat(path)
	if err != nil {
		logger.Panicln(err)
	}
	if !file.IsDir() {
		logger.Fatalln("Error: Target is not a directory.")
	}
}

func main() {
	flag.Parse()

	logger := log.New(os.Stdout, "", 0)
	if *showVersion {
		logger.Println(version)
		return
	}

	validateDir(mediahandler.MediaDir, logger)
	if *output {
		if err := mediahandler.WriteToStdout(); err != nil {
			logger.Panicln(err)
		}
		return
	}

	var handler http.Handler
	logger.SetFlags(log.Ldate | log.Ltime)
	if *password != "" {
		handler = router.NewAuthWrapper(router.InitRouter(logger), logger, *username, *password)
	} else {
		handler = router.NewLogWrapper(router.InitRouter(logger), logger)
	}

	addr := net.JoinHostPort(*listenAddr, strconv.Itoa(*listenPort))
	if *certFile != "" && *keyFile != "" {
		logger.Printf("Starting xplay server %s at https://%s/ ...\n", version, addr)
		if err := http.ListenAndServeTLS(addr, *certFile, *keyFile, handler); err != nil {
			logger.Panicln(err)
		}
	} else {
		logger.Printf("Starting xplay server %s at http://%s/ ...\n", version, addr)
		if err := http.ListenAndServe(addr, handler); err != nil {
			logger.Panicln(err)
		}
	}
}
