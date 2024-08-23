package router

import (
	"net/http"

	"octane.top/xplay/mediahandler"
	"octane.top/xplay/xspf"
)

const (
	xspfPath      = "/play.xspf"
	mediaBasePath = "/media/"
	imageBasePath = "/img/"
)

func httpHandler(w http.ResponseWriter, _ *http.Request) {
	playList, err := mediahandler.GetMedia(mediaBasePath, imageBasePath)
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
