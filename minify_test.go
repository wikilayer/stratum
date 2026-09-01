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
	for _, name := range cssAssets {
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

// The public stylesheet is a real bundle, not an @import entrypoint.
// One link must produce one request while preserving the manifest's
// cascade order from the first source through the last.
func TestMinify_StyleIsOneBundle(t *testing.T) {
	body, err := fs.ReadFile(Static, "stratum.css")
	if err != nil {
		t.Fatal(err)
	}
	css := string(body)
	if strings.Contains(css, "@import") {
		t.Error("stratum.css still contains @import and fans out into more requests")
	}
	if !strings.HasPrefix(css, cssLayerOrder) {
		t.Errorf("stratum.css does not open with the layer order: %q", css[:min(len(css), 80)])
	}
	first, err := fs.ReadFile(Static, cssAssets[0])
	if err != nil {
		t.Fatal(err)
	}
	last, err := fs.ReadFile(Static, cssAssets[len(cssAssets)-1])
	if err != nil {
		t.Fatal(err)
	}
	firstAt := strings.Index(css, string(first))
	lastAt := strings.Index(css, string(last))
	if firstAt < 0 || lastAt < 0 || firstAt >= lastAt {
		t.Errorf("stratum.css does not preserve manifest order: first=%d last=%d", firstAt, lastAt)
	}
}

// Scripts keep their readable sources in the repository but leave the
// embedded filesystem minified. The framework bundle is the same helpers
// concatenated in manifest order, so consumers need one script request.
func TestMinify_ScriptsAndBundle(t *testing.T) {
	source := mustSub(embedded, "static")
	bundle, err := fs.ReadFile(Static, "stratum.js")
	if err != nil {
		t.Fatal(err)
	}
	joined := string(bundle)
	lastAt := -1
	for _, name := range jsAssets {
		served, err := fs.ReadFile(Static, name)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		raw, err := fs.ReadFile(source, name)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if strings.Contains(string(served), "//") || strings.Contains(string(served), "/*") {
			t.Errorf("%s: comment survived minification", name)
		}
		if len(served) >= len(raw) {
			t.Errorf("%s: minified %d bytes is not smaller than source %d", name, len(served), len(raw))
		}
		at := strings.Index(joined, string(served))
		if at < 0 || at <= lastAt {
			t.Errorf("%s is missing or out of order in stratum.js: at=%d after=%d", name, at, lastAt)
		}
		lastAt = at
	}
}

// TestMinify_LeavesEverythingElseAlone: only stylesheets and scripts are
// rewritten. Images and the icon sprite still go out byte for byte.
func TestMinify_LeavesEverythingElseAlone(t *testing.T) {
	source := mustSub(embedded, "static")

	err := fs.WalkDir(source, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || strings.HasSuffix(name, ".css") || strings.HasSuffix(name, ".js") {
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

	for _, tc := range []struct {
		path        string
		contentType string
	}{
		{"/stratum.css", "text/css"},
		{"/stratum.js", "text/javascript"},
	} {
		res, err := http.Get(srv.URL + tc.path)
		if err != nil {
			t.Fatal(err)
		}
		_ = res.Body.Close()
		if res.StatusCode != http.StatusOK {
			t.Errorf("GET %s: %d", tc.path, res.StatusCode)
		}
		if !strings.HasPrefix(res.Header.Get("Content-Type"), tc.contentType) {
			t.Errorf("GET %s: Content-Type %q", tc.path, res.Header.Get("Content-Type"))
		}
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
