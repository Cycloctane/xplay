package router

import (
	"fmt"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/Cycloctane/xplay/pkg/xspf"
)

const testHostname = "yes.local:8080"

func testResponse(t *testing.T, tmpDir, scheme, expectedUrl string) {
	logger := log.New(io.Discard, "", 0)
	handler := InitRouter(tmpDir, scheme, logger)
	req := httptest.NewRequest("GET", xspfPath, nil)
	req.Host = testHostname
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	if rec.Header().Get("Content-Type") != xspf.ContentType {
		t.Error("Error Content-Type in response.")
	}
	list, err := xspf.DecodeXspf(rec.Body)
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Tracks) < 1 {
		t.Error("Expected at least one track in the response.")
	}
	if list.Tracks[0].Location != expectedUrl {
		t.Errorf("Expected location URL %s, got %s", expectedUrl, list.Tracks[0].Location)
	}
}

func TestResponse(t *testing.T) {
	tmpDir := t.TempDir()
	f, err := os.Create(path.Join(tmpDir, "test.mkv"))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Run("relative path", func(t *testing.T) {
		testResponse(t, tmpDir, "", "/media/test.mkv")
	})
	t.Run("absolute path", func(t *testing.T) {
		testResponse(t, tmpDir, "http", fmt.Sprintf("http://%s/media/test.mkv", testHostname))
		testResponse(t, tmpDir, "https", fmt.Sprintf("https://%s/media/test.mkv", testHostname))
	})
}
