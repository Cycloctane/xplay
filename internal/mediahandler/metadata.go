package mediahandler

import (
	"io"
	"os"
	"strconv"

	"github.com/Cycloctane/xplay/pkg/xspf"
	"github.com/dhowden/tag"
)

// Keys: media file extensions. Values: if metadata parsing of this ext is supported or not
var supportedExt = map[string]bool{".mp3": true, ".flac": true, ".ogg": true, ".mp4": true, ".mkv": false}

func addTag(t *string, m string) {
	if m != "" {
		*t = m
	}
}

func readTag(f *os.File, track *xspf.Track) error {
	metadata, err := tag.ReadFrom(f)
	if err != nil {
		return err
	}
	addTag(&track.Title, metadata.Title())
	addTag(&track.Creator, metadata.Artist())
	addTag(&track.Album, metadata.Album())
	if trackNum, _ := metadata.Track(); trackNum > 0 {
		addTag(&track.TrackNum, strconv.Itoa(trackNum))
	}
	addTag(&track.Annotation, metadata.Comment())
	if metadata.Picture() != nil {
		addTag(&track.ImageExt, metadata.Picture().Ext)
	}
	return nil
}

func readImg(f io.ReadSeeker) (*tag.Picture, error) {
	metadata, err := tag.ReadFrom(f)
	if err != nil {
		return nil, err
	}
	return metadata.Picture(), nil
}
