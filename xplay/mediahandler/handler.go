package mediahandler

import (
	"bufio"
	"fmt"
	"net/url"
	"octane.top/xplay/xspf"
	"os"
	"path/filepath"
	"strings"
)

var MediaDir string

var taggedExt = [...]string{".mp3", ".flac", ".ogg", "mp4"}

func isTaggedExt(ext string) bool {
	for _, v := range taggedExt {
		if ext == v {
			return true
		}
	}
	return false
}

func getTaggedTrack(path string, track *xspf.Track) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := ReadTag(f, track); err != nil {
		return err
	}
	return nil
}

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
		ext := filepath.Ext(v.Name())
		track := &xspf.Track{
			Location: location,
			Title:    strings.TrimSuffix(v.Name(), ext),
		}
		if isTaggedExt(ext) {
			path := filepath.Join(MediaDir, v.Name())
			if err := getTaggedTrack(path, track); err != nil {
				fmt.Printf("cannot parse %s: %s\n", path, err.Error())
				continue
			}
		}
		playList.TrackList.Tracks = append(playList.TrackList.Tracks, *track)
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
