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
	mediaDir    = flag.String("d", ".", "served directory")
	output      = flag.Bool("w", false, "write xspf to stdout and exit")
	absolute    = flag.Bool("absolute-links", false, "use absolute links in playlist (compatible with more media players)")
	listenPort  = flag.Int("p", defaultPort, "http server bind port")
	listenAddr  = flag.String("b", "0.0.0.0", "http server bind address")
	username    = flag.String("username", defaultUser, "http basic auth username")
	password    = flag.String("password", "", "http basic auth password")
	certFile    = flag.String("ssl-cert", "", "cert file path for https support")
	keyFile     = flag.String("ssl-key", "", "cert key path for https support")
	hostnames   = flag.String("server-hostname", "", "comma separated server hostnames for Host header validation")
)

func init() {
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

	validateDir(*mediaDir, logger)
	if *output {
		if err := mediahandler.WriteToStdout(*mediaDir); err != nil {
			logger.Panicln(err)
		}
		return
	}

	var scheme string
	if *certFile != "" && *keyFile != "" {
		scheme = "https"
	} else {
		scheme = "http"
	}

	var handler http.Handler
	logger.SetFlags(log.Ldate | log.Ltime)
	if *absolute {
		handler = router.InitRouter(*mediaDir, scheme, logger)
	} else {
		handler = router.InitRouter(*mediaDir, "", logger)
	}

	if *password != "" {
		handler = router.NewBasicAuth(handler, logger, *username, *password)
	}
	if *hostnames != "" {
		handler = router.NewHostHeaderValidator(handler, *hostnames)
	}
	handler = router.NewLogWrapper(handler, logger)

	addr := net.JoinHostPort(*listenAddr, strconv.Itoa(*listenPort))
	logger.Printf("Starting xplay server %s at %s://%s/ ...\n", version, scheme, addr)
	if scheme == "https" {
		if err := http.ListenAndServeTLS(addr, *certFile, *keyFile, handler); err != nil {
			logger.Panicln(err)
		}
	} else if err := http.ListenAndServe(addr, handler); err != nil {
		logger.Panicln(err)
	}
}
