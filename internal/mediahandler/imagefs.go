package mediahandler

import (
	"bytes"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"time"
)

type mediaImageInfo struct {
	OriginFileInfo fs.FileInfo
	ImageName      string
	Length         int
}

func (mi mediaImageInfo) IsDir() bool {
	return false
}

func (mi mediaImageInfo) Sys() any {
	return nil
}

func (mi mediaImageInfo) Name() string {
	return mi.ImageName
}

func (mi mediaImageInfo) Size() int64 {
	return int64(mi.Length)
}

func (mi mediaImageInfo) Mode() fs.FileMode {
	return mi.OriginFileInfo.Mode()
}

func (mi mediaImageInfo) ModTime() time.Time {
	return mi.OriginFileInfo.ModTime()
}

type mediaImage struct {
	OriginFile fs.File
	Info       *mediaImageInfo
	*bytes.Reader
}

func (m mediaImage) Close() error {
	return m.OriginFile.Close()
}

func (m mediaImage) Readdir(_ int) ([]fs.FileInfo, error) {
	return []fs.FileInfo{}, nil
}

func (m mediaImage) Stat() (fs.FileInfo, error) {
	return m.Info, nil
}

type ImageFS struct {
	Mfs *MediaFS
}

func (ifs *ImageFS) Open(name string) (http.File, error) {
	f, err := ifs.Mfs.Open(name)
	if err != nil {
		return nil, err
	}
	pic, err := ReadImg(f)
	if err != nil {
		return nil, err
	}
	if pic == nil {
		return nil, os.ErrNotExist
	}
	fInfo, err := f.Stat()
	if err != nil {
		return nil, err
	}
	imgBytes := bytes.NewReader(pic.Data)
	imgInfo := &mediaImageInfo{
		OriginFileInfo: fInfo,
		ImageName:      fmt.Sprintf("album.%s", pic.Ext),
		Length:         imgBytes.Len(),
	}
	imgFile := mediaImage{OriginFile: f, Info: imgInfo, Reader: imgBytes}
	return imgFile, nil
}
