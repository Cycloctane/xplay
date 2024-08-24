package mediahandler

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

var supportedExt = [...]string{".mp3", ".flac", ".ogg", ".mp4", ".mkv"}

var Recursive bool

func validateFileType(f fs.DirEntry) bool {
	if f.IsDir() || f.Name()[0] == '.' {
		return false
	}
	for _, v := range supportedExt {
		if filepath.Ext(f.Name()) == v {
			return true
		}
	}
	return false
}

type MediaFS struct {
	Fs http.FileSystem
}

func (mfs *MediaFS) Open(name string) (http.File, error) {
	if !Recursive && filepath.Dir(name) != string(os.PathSeparator) {
		return nil, os.ErrNotExist
	}
	f, err := mfs.Fs.Open(name)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !validateFileType(fs.FileInfoToDirEntry(info)) {
		return nil, os.ErrNotExist
	}
	return f, nil
}
