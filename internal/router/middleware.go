package router

import (
	"crypto/sha256"
	"crypto/subtle"
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

type BasicAuth struct {
	http.Handler
	logger         *log.Logger
	username       string
	passwordSha224 [28]byte
}

func (aw *BasicAuth) verify(inputUsername, inputPassword string) bool {
	if inputUsername != aw.username {
		return false
	}
	inputPasswordHash := sha256.Sum224([]byte(inputPassword))
	return subtle.ConstantTimeCompare(aw.passwordSha224[:], inputPasswordHash[:]) == 1
}

func unauthorizedHandler(w http.ResponseWriter) {
	w.Header().Add("WWW-Authenticate", "Basic")
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func (aw *BasicAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw := &wrappedResponseWriter{ResponseWriter: w}
	user, password, isOk := r.BasicAuth()
	if !isOk {
		unauthorizedHandler(rw)
	} else if !aw.verify(user, password) {
		unauthorizedHandler(rw)
		aw.logger.Printf("[auth] %s - Authorization Failed %s:%s", r.RemoteAddr, user, password)
	} else {
		aw.Handler.ServeHTTP(rw, r)
	}
}

func NewBasicAuth(handler http.Handler, logger *log.Logger, username, password string) *BasicAuth {
	newHandler := &BasicAuth{
		Handler: handler, logger: logger, username: username,
	}
	newHandler.passwordSha224 = sha256.Sum224([]byte(password))
	return newHandler
}

type HostHeaderValidator struct {
	http.Handler
	hostnames map[string]struct{}
}

func (hv *HostHeaderValidator) verifyHost(host string) bool {
	host = strings.ToLower(host)
	if hostname, _, _ := net.SplitHostPort(host); hostname != "" {
		host = hostname
	}
	_, have := hv.hostnames[host]
	return have
}

func (hv *HostHeaderValidator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !hv.verifyHost(r.Host) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	hv.Handler.ServeHTTP(w, r)
}

func NewHostHeaderValidator(handler http.Handler, trustedHostHeaders string) *HostHeaderValidator {
	hostnames := make(map[string]struct{})
	for host := range strings.SplitSeq(trustedHostHeaders, ",") {
		hostnames[host] = struct{}{}
	}
	return &HostHeaderValidator{Handler: handler, hostnames: hostnames}
}
