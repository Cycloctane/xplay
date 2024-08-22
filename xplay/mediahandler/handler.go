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

const fileBaseURL = "file:///"

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

func GetMedia(MediaBase string, ImageBase string) (*xspf.PlayList, error) {
	MediaBaseURL, err := url.Parse(MediaBase)
	if err != nil {
		return nil, err
	}
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
		location := MediaBaseURL.JoinPath(v.Name())
		ext := filepath.Ext(v.Name())
		track := &xspf.Track{
			Location: location.String(),
			Title:    strings.TrimSuffix(v.Name(), ext),
		}
		if isTaggedExt(ext) {
			mediaFilePath := filepath.Join(MediaDir, v.Name())
			if err := getTaggedTrack(mediaFilePath, track); err != nil {
				fmt.Printf("cannot parse %s: %s\n", mediaFilePath, err.Error())
				continue
			}
			if MediaBaseURL.Scheme != "file" && track.ImageExt != "" {
				track.ImageURI, _ = url.JoinPath(ImageBase, v.Name())
			}
		}
		playList.TrackList.Tracks = append(playList.TrackList.Tracks, *track)
	}
	return playList, nil
}

func WriteToStdout() error {
	baseURL, err := url.Parse(fileBaseURL)
	if err != nil {
		return err
	}
	absPath, err := filepath.Abs(MediaDir)
	if err != nil {
		return err
	}
	fileUrl := baseURL.JoinPath(filepath.ToSlash(absPath))
	playList, err := GetMedia(fileUrl.String(), "")
	if err != nil {
		return err
	}
	if err := xspf.Generate(bufio.NewWriter(os.Stdout), playList); err != nil {
		return err
	}
	return nil
}
