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

func InitRouter(mediaDir, urlScheme string, logger *log.Logger) *http.ServeMux {
	var getBaseUrl func(BasePath, rhost string) *url.URL
	if urlScheme == "" { // Use relative urls if scheme is not set
		getBaseUrl = func(basePath, _ string) *url.URL {
			baseUrl, _ := url.Parse(basePath)
			return baseUrl
		}
	} else { // absolute urls
		getBaseUrl = func(basePath, rhost string) *url.URL {
			baseUrl := &url.URL{Scheme: urlScheme, Host: rhost}
			return baseUrl.JoinPath(basePath)
		}
	}

	httpFS := &mediahandler.MediaFS{Fs: http.Dir(mediaDir)}
	router := http.NewServeMux()
	router.HandleFunc(xspfPath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		playList, err := mediahandler.GetMedia(
			mediaDir, getBaseUrl(mediaBasePath, r.Host), getBaseUrl(imageBasePath, r.Host),
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
		w.WriteHeader(http.StatusOK)
		buffered.WriteTo(w)
	})
	router.Handle(mediaBasePath, http.StripPrefix(mediaBasePath, http.FileServer(httpFS)))
	router.Handle(imageBasePath, http.StripPrefix(imageBasePath, http.FileServer(&mediahandler.ImageFS{Mfs: httpFS})))
	router.Handle("/favicon.ico", http.NotFoundHandler())
	router.Handle("/", http.RedirectHandler(xspfPath, http.StatusTemporaryRedirect))
	return router
}
