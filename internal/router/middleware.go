package router

import (
	"log"
	"net"
	"net/http"
	"strings"
)

const serverHeader = "xplay"

type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (wr *wrappedResponseWriter) WriteHeader(statusCode int) {
	wr.Header().Set("Server", serverHeader)
	wr.statusCode = statusCode
	wr.ResponseWriter.WriteHeader(statusCode)
}

func (wr *wrappedResponseWriter) Write(data []byte) (int, error) {
	s, err := wr.ResponseWriter.Write(data)
	wr.size += s
	return s, err
}

type LogWrapper struct {
	http.Handler
	logger *log.Logger
}

func (lw *LogWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw := &wrappedResponseWriter{ResponseWriter: w}
	lw.Handler.ServeHTTP(rw, r)
	lw.logger.Printf("[http] %s - \"%s %s\" %d %d \"%s\"", r.RemoteAddr, r.Method, r.RequestURI, rw.statusCode, rw.size, r.UserAgent())
}

func NewLogWrapper(handler http.Handler, logger *log.Logger) *LogWrapper {
	return &LogWrapper{handler, logger}
}

type TrustedHostWrapper struct {
	http.Handler
	trustedHosts []string
}

func (th *TrustedHostWrapper) verifyHost(host string) bool {
	for _, v := range th.trustedHosts {
		if host == v {
			return true
		}
		if h, _, _ := net.SplitHostPort(host); h == v {
			return true
		}
	}
	return false
}

func (th *TrustedHostWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !th.verifyHost(r.Host) {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return
	}
	th.Handler.ServeHTTP(w, r)
}

func NewTrustedHostWrapper(handler http.Handler, trustedHosts string) *TrustedHostWrapper {
	return &TrustedHostWrapper{Handler: handler, trustedHosts: strings.Split(trustedHosts, ",")}
}
