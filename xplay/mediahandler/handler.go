package mediahandler

import (
	"bufio"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"octane.top/xplay/xspf"
)

var MediaDir string

func GetMedia(baseURL string) (*xspf.PlayList, error) {
	playList := &xspf.PlayList{
		Version: "1", XMLns: "http://xspf.org/ns/0/",
		Creator: "xplay",
	}
	if err := fs.WalkDir(os.DirFS(MediaDir), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || !validateFileType(d) {
			return nil
		}
		if !Recursive && filepath.Dir(path) != "." {
			return nil
		}
		location, err := url.JoinPath(baseURL, path)
		if err != nil {
			return err
		}
		playList.TrackList.Tracks = append(playList.TrackList.Tracks, xspf.Track{
			Location: location,
			Title:    strings.TrimSuffix(d.Name(), filepath.Ext(d.Name())),
		})
		return nil
	}); err != nil {
		return nil, err
	}
	return playList, nil
}

func WriteToStdout() error {
	absPath, err := filepath.Abs(MediaDir)
	if err != nil {
		return err
	}
	fileUrl, err := url.JoinPath("file:///", filepath.ToSlash(absPath))
	if err != nil {
		return err
	}
	playList, err := GetMedia(fileUrl)
	if err != nil {
		return err
	}
	if err := xspf.Generate(bufio.NewWriter(os.Stdout), playList); err != nil {
		return err
	}
	return nil
}
