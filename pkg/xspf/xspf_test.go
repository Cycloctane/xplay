package xspf

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"testing"
)

const testXml = `
<playlist version="1" xmlns="http://xspf.org/ns/0/">
  <title>Playlist</title>
  <date>2005-01-08T17:10:47-05:00</date>
  <trackList>
    <track>
      <location>1.mp3</location>
      <title>1</title>
    </track>
    <track>
      <location>2.mp3</location>
      <title>2</title>
    </track>
  </trackList>
</playlist>
`

func newTestPlaylist() *PlayList {
	list := &PlayList{
		Title: "Playlist",
		Date:  "2005-01-08T17:10:47-05:00",
	}
	for i := 1; i < 3; i++ {
		list.Tracks = append(list.Tracks, &Track{
			Location: strconv.Itoa(i) + ".mp3", Title: strconv.Itoa(i), ImageExt: "jpg",
		})
	}
	return list
}

func TestEncodeXspf(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	if err := EncodeXspf(buf, newTestPlaylist()); err != nil {
		t.Fatal(err)
	}
	if strings.Trim(buf.String(), "\n") != strings.Trim(testXml, "\n") {
		t.Errorf("Encoded playlist does not match expected:\n%s", buf.String())
	}
}

func TestGenerate(t *testing.T) {
	if err := Generate(io.Discard, newTestPlaylist()); err != nil {
		t.Error(err)
	}
}

func TestBufferedGenerate(t *testing.T) {
	if _, err := BufferedGenerate(newTestPlaylist()); err != nil {
		t.Error(err)
	}
}
