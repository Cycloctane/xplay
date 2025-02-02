package mediahandler

import (
	"bufio"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Cycloctane/xplay/pkg/xspf"
)

const fileBaseURL = "file:///"

func getTaggedTrack(path string, track *xspf.Track) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := readTag(f, track); err != nil {
		return err
	}
	return nil
}

func GetMedia(MediaBaseURL, ImageBaseURL *url.URL) (*xspf.PlayList, error) {
	playList := &xspf.PlayList{Creator: "xplay", Title: "xplay"}
	if err := fs.WalkDir(os.DirFS(MediaDir), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !isSupportedFile(d) || (NoRecursive && filepath.Dir(path) != ".") {
			return nil
		}
		location := MediaBaseURL.JoinPath(path)
		ext := filepath.Ext(d.Name())
		track := &xspf.Track{
			Location: location.String(),
			Title:    strings.TrimSuffix(d.Name(), ext),
		}
		if !NoTag && supportedExt[ext] {
			mediaFilePath := filepath.Join(MediaDir, path)
			if err := getTaggedTrack(mediaFilePath, track); err != nil {
				return err
			}
			if MediaBaseURL.Scheme != "file" && track.ImageExt != "" {
				track.Image = ImageBaseURL.JoinPath(path).String()
			}
		}
		playList.Tracks = append(playList.Tracks, track)
		return nil
	}); err != nil {
		return nil, err
	}
	return playList, nil
}

func WriteToStdout() error {
	baseUrl, _ := url.Parse(fileBaseURL)
	EmptyUrl, _ := url.Parse("")
	absPath, err := filepath.Abs(MediaDir)
	if err != nil {
		return err
	}
	fileUrl := baseUrl.JoinPath(filepath.ToSlash(absPath))
	playList, err := GetMedia(fileUrl, EmptyUrl)
	if err != nil {
		return err
	}
	if err := xspf.Generate(bufio.NewWriter(os.Stdout), playList); err != nil {
		return err
	}
	return nil
}
