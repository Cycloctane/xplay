package mediahandler

import (
	"net/http"
	"os"
	"path/filepath"
)

var supportedExt = [...]string{".mp3", ".flac", ".ogg", ".mp4", ".mkv"}

type fileInfoOrDirEntry interface {
	Name() string
	IsDir() bool
}

func validateFileType(f fileInfoOrDirEntry) bool {
	ext := filepath.Ext(f.Name())
	if f.IsDir() {
		return false
	}
	for _, v := range supportedExt {
		if ext == v {
			return true
		}
	}
	return false
}

type MediaFS struct {
	Fs http.FileSystem
}

func (mfs MediaFS) Open(name string) (http.File, error) {
	f, err := mfs.Fs.Open(name)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !validateFileType(info) {
		return nil, os.ErrNotExist
	}
	return f, nil
}
