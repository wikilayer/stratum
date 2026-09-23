package main

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
	"github.com/wikilayer/stratum"
	"golang.org/x/net/html"
)

func TestHandlerServesTemplatesAndBundledAssetsWithRevalidation(t *testing.T) {
	handler := newHandler(testDesignSystem())

	page := request(t, handler, "/design-system/")
	require.Equal(t, http.StatusOK, page.Code)
	require.Equal(t, "no-cache", page.Header().Get("Cache-Control"))
	require.NotEmpty(t, page.Header().Get("ETag"))
	document, err := html.Parse(page.Body)
	require.NoError(t, err)
	require.Equal(t, "/static/stratum.css", stylesheetURL(document))

	cachedPage := requestWithETag(t, handler, "/design-system/", page.Header().Get("ETag"))
	require.Equal(t, http.StatusNotModified, cachedPage.Code)
	require.Empty(t, cachedPage.Body.Bytes())

	asset := request(t, handler, "/static/stratum.css")
	require.Equal(t, http.StatusOK, asset.Code)
	require.Equal(t, "public, max-age=0, must-revalidate", asset.Header().Get("Cache-Control"))
	wantAsset, err := fs.ReadFile(stratum.Static, "stratum.css")
	require.NoError(t, err)
	require.Equal(t, wantAsset, asset.Body.Bytes())
	require.NotEmpty(t, asset.Header().Get("ETag"))

	cachedAsset := requestWithETag(t, handler, "/static/stratum.css", asset.Header().Get("ETag"))
	require.Equal(t, http.StatusNotModified, cachedAsset.Code)
	require.Empty(t, cachedAsset.Body.Bytes())

	partial := request(t, handler, "/design-system/_document.gohtml")
	require.Equal(t, http.StatusNotFound, partial.Code)
}

func TestBuildDiscoversPagesAndCopiesAssets(t *testing.T) {
	destination := filepath.Join(t.TempDir(), "site")
	require.NoError(t, build(testDesignSystem(), destination))

	page, err := os.Open(filepath.Join(destination, "design-system", "new-example.html"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, page.Close()) })
	document, err := html.Parse(page)
	require.NoError(t, err)
	require.Equal(t, "../static/stratum.css", stylesheetURL(document))

	image, err := os.ReadFile(filepath.Join(destination, "design-system", "example.svg"))
	require.NoError(t, err)
	require.Equal(t, []byte("<svg/>"), image)

	_, err = os.Stat(filepath.Join(destination, "design-system", sharedLayout))
	require.ErrorIs(t, err, fs.ErrNotExist)

	gotAsset, err := os.ReadFile(filepath.Join(destination, "static", "stratum.css"))
	require.NoError(t, err)
	wantAsset, err := fs.ReadFile(stratum.Static, "stratum.css")
	require.NoError(t, err)
	require.Equal(t, wantAsset, gotAsset)
}

func testDesignSystem() fstest.MapFS {
	return fstest.MapFS{
		sharedLayout: {
			Data: []byte(`{{define "document-head"}}<!doctype html><html><head><link rel="stylesheet" href="{{asset "stratum.css"}}"></head><body>{{end}}{{define "document-end"}}</body></html>{{end}}`),
		},
		"index.html": {
			Data: []byte(`{{template "document-head" .}}<main>Index</main>{{template "document-end"}}`),
		},
		"new-example.html": {
			Data: []byte(`{{template "document-head" .}}<main>Automatically discovered</main>{{template "document-end"}}`),
		},
		"example.svg": {Data: []byte("<svg/>")},
	}
}

func request(t *testing.T, handler http.Handler, target string) *httptest.ResponseRecorder {
	t.Helper()
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, target, nil))
	return recorder
}

func requestWithETag(t *testing.T, handler http.Handler, target, etag string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodGet, target, nil)
	request.Header.Set("If-None-Match", etag)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}

func stylesheetURL(document *html.Node) string {
	var visit func(*html.Node) string
	visit = func(node *html.Node) string {
		if node.Type == html.ElementNode && node.Data == "link" && attribute(node, "rel") == "stylesheet" {
			return attribute(node, "href")
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if href := visit(child); href != "" {
				return href
			}
		}
		return ""
	}
	return visit(document)
}

func attribute(node *html.Node, name string) string {
	for _, attribute := range node.Attr {
		if attribute.Key == name {
			return attribute.Val
		}
	}
	return ""
}
