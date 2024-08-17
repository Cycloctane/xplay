package xspf

import (
	"bytes"
	"encoding/xml"
	"io"
)

const ContentType = "application/xspf+xml"

type PlayList struct {
	XMLName    xml.Name `xml:"playlist"`
	Version    string   `xml:"version,attr"`
	XMLns      string   `xml:"xmlns,attr"`
	Title      string   `xml:"title,omitempty"`
	Creator    string   `xml:"creator,omitempty"`
	Date       string   `xml:"date,omitempty"`
	Annotation string   `xml:"annotation,omitempty"`
	TrackList  TrackList
}

type TrackList struct {
	XMLName xml.Name `xml:"trackList"`
	Tracks  []Track
}

type Track struct {
	XMLName    xml.Name `xml:"track"`
	Location   string   `xml:"location"`
	Title      string   `xml:"title,omitempty"`
	Creator    string   `xml:"creator,omitempty"`
	Album      string   `xml:"album,omitempty"`
	TrackNum   string   `xml:"trackNum,omitempty"`
	Duration   string   `xml:"duration,omitempty"`
	Annotation string   `xml:"annotation,omitempty"`
}

func Generate(w io.Writer, list *PlayList) error {
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")
	if err := encoder.Encode(list); err != nil {
		return err
	}
	return nil
}

func BufferedGenerate(list *PlayList) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer([]byte(xml.Header))
	encoder := xml.NewEncoder(buf)
	encoder.Indent("", "  ")
	if err := encoder.Encode(list); err != nil {
		return nil, err
	}
	return buf, nil
}
