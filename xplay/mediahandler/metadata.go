package mediahandler

import (
	"github.com/dhowden/tag"
	"io"
	"octane.top/xplay/xspf"
	"os"
	"strconv"
)

func addTag(t *string, m string) {
	if m != "" {
		*t = m
	}
}

func ReadTag(f *os.File, track *xspf.Track) error {
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

func ReadImg(f io.ReadSeeker) (*tag.Picture, error) {
	metadata, err := tag.ReadFrom(f)
	if err != nil {
		return nil, err
	}
	return metadata.Picture(), nil
}
