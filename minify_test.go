package stratum

import (
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// TestMinify_StylesheetsLoseTheirProse: the served CSS carries no
// comments and is materially smaller than the source, which is the
// whole point; and it still contains the rules, so "smaller" was not
// achieved by serving nothing.
func TestMinify_StylesheetsLoseTheirProse(t *testing.T) {
	source := mustSub(embedded, "static")

	var served, raw int
	for _, name := range CSSAssets {
		body, err := fs.ReadFile(Static, name)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		original, err := fs.ReadFile(source, name)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		served += len(body)
		raw += len(original)

		if strings.Contains(string(body), "/*") {
			t.Errorf("%s: comment survived minification", name)
		}
		if len(body) >= len(original) {
			t.Errorf("%s: minified %d bytes is not smaller than source %d", name, len(body), len(original))
		}
	}

	if served*2 > raw {
		t.Errorf("minified bundle is %d bytes against %d of source: expected to at least halve", served, raw)
	}
}

// TestMinify_LeavesEverythingElseAlone: only stylesheets are rewritten.
// The icon sprite and the scripts go out byte for byte, so a sprite id
// or a script's behaviour cannot change under the consumer's feet.
func TestMinify_LeavesEverythingElseAlone(t *testing.T) {
	source := mustSub(embedded, "static")

	err := fs.WalkDir(source, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || strings.HasSuffix(name, ".css") {
			return nil
		}
		want, err := fs.ReadFile(source, name)
		if err != nil {
			return err
		}
		got, err := fs.ReadFile(Static, name)
		if err != nil {
			return err
		}
		if string(got) != string(want) {
			t.Errorf("%s: served bytes differ from the source", name)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// TestMinify_ServesOverHTTP: the tree still answers what a file server
// asks of it. A stylesheet's Content-Length has to be the length of
// what is actually sent, or the browser waits for bytes that never
// arrive; directories still list; a missing name is still a 404.
func TestMinify_ServesOverHTTP(t *testing.T) {
	srv := httptest.NewServer(http.FileServer(http.FS(Static)))
	defer srv.Close()

	res, err := http.Get(srv.URL + "/css/utilities.css")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("GET css/utilities.css: %d", res.StatusCode)
	}
	body := make([]byte, res.ContentLength+1)
	n, _ := io.ReadFull(res.Body, body)
	if int64(n) != res.ContentLength {
		t.Errorf("Content-Length says %d, body is %d bytes", res.ContentLength, n)
	}
	if got := res.Header.Get("Content-Length"); got != strconv.Itoa(n) {
		t.Errorf("Content-Length header %q against %d bytes sent", got, n)
	}

	missing, err := http.Get(srv.URL + "/css/nothing-here.css")
	if err != nil {
		t.Fatal(err)
	}
	defer missing.Body.Close()
	if missing.StatusCode != http.StatusNotFound {
		t.Errorf("missing stylesheet: got %d, want 404", missing.StatusCode)
	}
}
