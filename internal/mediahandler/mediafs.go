package mediahandler

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

func isSupportedFile(f fs.DirEntry) bool {
	if f.IsDir() || f.Name()[0] == '.' {
		return false
	}
	_, exists := supportedExt[filepath.Ext(f.Name())]
	return exists
}

type MediaFS struct {
	Fs http.FileSystem
}

func (mfs *MediaFS) Open(name string) (http.File, error) {
	if NoRecursive && filepath.Dir(name) != string(os.PathSeparator) {
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
	if !isSupportedFile(fs.FileInfoToDirEntry(info)) {
		return nil, os.ErrNotExist
	}
	return f, nil
}
