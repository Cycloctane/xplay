package mediahandler

import (
	"bufio"
	"net/url"
	"octane.top/xplay/xspf"
	"os"
	"path/filepath"
	"strings"
)

var MediaDir string

func GetMedia(baseURL string) (*xspf.PlayList, error) {
	files, err := os.ReadDir(MediaDir)
	if err != nil {
		return nil, err
	}
	playList := &xspf.PlayList{
		Version: "1", XMLns: "http://xspf.org/ns/0/",
		Creator: "xplay",
	}
	playList.TrackList.Tracks = make([]xspf.Track, 0, len(files))
	for _, v := range files {
		if !validateFileType(v) {
			continue
		}
		location, _ := url.JoinPath(baseURL, url.PathEscape(v.Name()))
		playList.TrackList.Tracks = append(playList.TrackList.Tracks, xspf.Track{
			Location: location,
			Title:    strings.TrimSuffix(v.Name(), filepath.Ext(v.Name())),
		})
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
