package router

import (
	"crypto/sha256"
	"crypto/subtle"
	"log"
	"net/http"
)

type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (wr *wrappedResponseWriter) WriteHeader(statusCode int) {
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

type AuthWrapper struct {
	http.Handler
	logger         *log.Logger
	username       string
	passwordSha224 [28]byte
}

func (aw *AuthWrapper) verify(inputUsername, inputPassword string) bool {
	if inputUsername != aw.username {
		return false
	}
	inputPasswordHash := sha256.Sum224([]byte(inputPassword))
	return subtle.ConstantTimeCompare(aw.passwordSha224[:], inputPasswordHash[:]) == 1
}

func unauthorizedHandler(w http.ResponseWriter) {
	w.Header().Add("WWW-Authenticate", "Basic")
	http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
}

func (aw *AuthWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	aw.logger.Printf("[http] %s - \"%s %s\" %d %d \"%s\"", r.RemoteAddr, r.Method, r.RequestURI, rw.statusCode, rw.size, r.UserAgent())
}

func NewAuthWrapper(handler http.Handler, logger *log.Logger, username, password string) *AuthWrapper {
	newHandler := &AuthWrapper{
		Handler: handler, logger: logger, username: username,
	}
	newHandler.passwordSha224 = sha256.Sum224([]byte(password))
	return newHandler
}
