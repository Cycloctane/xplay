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
)

func httpHandlerFactory(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		mediaBaseUrl, _ := url.Parse(mediaBasePath)
		imageBaseUrl, _ := url.Parse(imageBasePath)
		playList, err := mediahandler.GetMedia(mediaBaseUrl, imageBaseUrl)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Printf("[file] Error parsing media file: %v\n", err)
			return
		}
		buffered, err := xspf.BufferedGenerate(playList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Printf("[xspf] Error generating playlist: %v\n", err)
			return
		}
		w.Header().Set("Content-Type", xspf.ContentType)
		w.WriteHeader(http.StatusOK)
		buffered.WriteTo(w)
	}
}

func InitRouter(logger *log.Logger) *http.ServeMux {
	httpFS := &mediahandler.MediaFS{Fs: http.Dir(mediahandler.MediaDir)}
	router := http.NewServeMux()
	httpHandler := httpHandlerFactory(logger)
	router.HandleFunc(xspfPath, httpHandler)
	router.Handle(mediaBasePath, http.StripPrefix(mediaBasePath, http.FileServer(httpFS)))
	router.Handle(imageBasePath, http.StripPrefix(imageBasePath, http.FileServer(&mediahandler.ImageFS{Mfs: httpFS})))
	router.Handle("/favicon.ico", http.NotFoundHandler())
	router.Handle("/", http.RedirectHandler(xspfPath, http.StatusTemporaryRedirect))
	return router
}
