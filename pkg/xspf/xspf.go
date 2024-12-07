package xspf

import (
	"bytes"
	"encoding/xml"
	"io"
)

const (
	XmlVersion  = "1"
	Xmlns       = "http://xspf.org/ns/0/"
	ContentType = "application/xspf+xml"
)

type PlayList struct {
	Title      string   `xml:"title,omitempty"`
	Creator    string   `xml:"creator,omitempty"`
	Date       string   `xml:"date,omitempty"`
	Annotation string   `xml:"annotation,omitempty"`
	Tracks     []*Track `xml:"trackList>track"`
}

type Track struct {
	Location   string `xml:"location"`
	Title      string `xml:"title,omitempty"`
	Creator    string `xml:"creator,omitempty"`
	Album      string `xml:"album,omitempty"`
	TrackNum   string `xml:"trackNum,omitempty"`
	Duration   string `xml:"duration,omitempty"`
	ImageExt   string `xml:"-"`
	Image      string `xml:"image,omitempty"`
	Annotation string `xml:"annotation,omitempty"`
}

func EncodeXspf(w io.Writer, list *PlayList) error {
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")
	head := xml.StartElement{
		Name: xml.Name{Local: "playlist"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "version"}, Value: XmlVersion},
			{Name: xml.Name{Local: "xmlns"}, Value: Xmlns},
		},
	}
	return encoder.EncodeElement(list, head)
}

func Generate(w io.Writer, list *PlayList) error {
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	if err := EncodeXspf(w, list); err != nil {
		return err
	}
	return nil
}

func BufferedGenerate(list *PlayList) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer([]byte(xml.Header))
	if err := EncodeXspf(buf, list); err != nil {
		return nil, err
	}
	return buf, nil
}
