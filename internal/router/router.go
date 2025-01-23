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

func httpHandlerFactory(scheme string, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		baseUrl := &url.URL{Scheme: scheme, Host: r.Host}
		playList, err := mediahandler.GetMedia(
			baseUrl.JoinPath(mediaBasePath), baseUrl.JoinPath(imageBasePath),
		)
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

		// Location links in xspf are generated from request's Host Header.
		// To mitigate potential cache poisoning, make sure untrusted xspf responses are never cached.
		w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0")

		w.WriteHeader(http.StatusOK)
		buffered.WriteTo(w)
	}
}

func InitRouter(scheme string, logger *log.Logger) *http.ServeMux {
	httpFS := &mediahandler.MediaFS{Fs: http.Dir(mediahandler.MediaDir)}
	router := http.NewServeMux()
	httpHandler := httpHandlerFactory(scheme, logger)
	router.HandleFunc(xspfPath, httpHandler)
	router.Handle(mediaBasePath, http.StripPrefix(mediaBasePath, http.FileServer(httpFS)))
	router.Handle(imageBasePath, http.StripPrefix(imageBasePath, http.FileServer(&mediahandler.ImageFS{Mfs: httpFS})))
	router.Handle("/favicon.ico", http.NotFoundHandler())
	router.Handle("/", http.RedirectHandler(xspfPath, http.StatusTemporaryRedirect))
	return router
}
