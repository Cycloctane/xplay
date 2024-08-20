package mediahandler

import (
	"github.com/dhowden/tag"
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
	return nil
}
