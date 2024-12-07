package router

import (
	"log"
	"net/http"
	"net/url"

	"github.com/Cycloctane/xplay/internal/mediahandler"
	"github.com/Cycloctane/xplay/pkg/xspf"
)

const (
	xspfPath      = "/play.xspf"
	mediaBasePath = "/media/"
	imageBasePath = "/img/"
	serverHeader  = "xplay"
)

func httpHandler(w http.ResponseWriter, _ *http.Request) {
	mediaBaseUrl, _ := url.Parse(mediaBasePath)
	imageBaseUrl, _ := url.Parse(imageBasePath)
	w.Header().Set("Server", serverHeader)
	playList, err := mediahandler.GetMedia(mediaBaseUrl, imageBaseUrl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	buffered, err := xspf.BufferedGenerate(playList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", xspf.ContentType)
	w.WriteHeader(http.StatusOK)
	if _, err = buffered.WriteTo(w); err != nil {
		return
	}
}

func InitRouter() *http.ServeMux {
	httpFS := &mediahandler.MediaFS{Fs: http.Dir(mediahandler.MediaDir)}
	router := http.NewServeMux()
	router.HandleFunc(xspfPath, httpHandler)
	router.Handle(mediaBasePath, http.StripPrefix(mediaBasePath, http.FileServer(httpFS)))
	router.Handle(imageBasePath, http.StripPrefix(imageBasePath, http.FileServer(&mediahandler.ImageFS{Mfs: httpFS})))
	router.Handle("/favicon.ico", http.NotFoundHandler())
	router.Handle("/", http.RedirectHandler(xspfPath, http.StatusTemporaryRedirect))
	return router
}

func InitLogRouter(logger *log.Logger) http.Handler {
	return NewLogWrapper(InitRouter(), logger)
}

func InitAuthRouter(logger *log.Logger, username, password string) http.Handler {
	return NewAuthWrapper(InitRouter(), logger, username, password)
}
