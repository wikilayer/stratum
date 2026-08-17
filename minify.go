package stratum

import (
	"bytes"
	"io"
	"io/fs"
	"path"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
)

// minifiedFS serves the CSS of an asset tree with its comments and
// whitespace removed, and everything else byte for byte.
//
// The prose in these stylesheets is written for whoever edits them and
// runs to roughly two thirds of every file; a browser pays for it on
// each cold load and reads none of it. Minifying happens once, when the
// package initialises, rather than in a build step: a generated copy in
// the repository is a second source of truth that a forgotten
// regeneration silently ships stale, and a consumer would have no way
// to tell.
type minifiedFS struct {
	source fs.FS
	css    map[string][]byte
}

func newMinifiedFS(source fs.FS) fs.FS {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)

	out := &minifiedFS{source: source, css: map[string][]byte{}}
	err := fs.WalkDir(source, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || path.Ext(name) != ".css" {
			return nil
		}
		raw, err := fs.ReadFile(source, name)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		if err := m.Minify("text/css", &buf, bytes.NewReader(raw)); err != nil {
			return err
		}
		out.css[name] = buf.Bytes()
		return nil
	})
	// A stylesheet the minifier rejects is a stylesheet no browser will
	// read either, so it fails the process that embeds it rather than
	// reaching a page.
	if err != nil {
		panic("stratum: minify assets: " + err.Error())
	}
	return out
}

func (m *minifiedFS) Open(name string) (fs.File, error) {
	body, ok := m.css[name]
	if !ok {
		return m.source.Open(name)
	}
	info, err := fs.Stat(m.source, name)
	if err != nil {
		return nil, err
	}
	return &minifiedFile{
		Reader: bytes.NewReader(body),
		info:   minifiedInfo{FileInfo: info, size: int64(len(body))},
	}, nil
}

// ReadDir and Stat pass straight through: http.FileServer reaches for
// both, and a directory listing or a modification time is the same
// whether or not the bytes behind it were minified. Stat of a
// stylesheet goes through Open so its size is the size actually served.
func (m *minifiedFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(m.source, name)
}

func (m *minifiedFS) Stat(name string) (fs.FileInfo, error) {
	if _, ok := m.css[name]; !ok {
		return fs.Stat(m.source, name)
	}
	f, err := m.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.Stat()
}

type minifiedFile struct {
	*bytes.Reader
	info minifiedInfo
}

func (f *minifiedFile) Stat() (fs.FileInfo, error) { return f.info, nil }
func (f *minifiedFile) Close() error               { return nil }

type minifiedInfo struct {
	fs.FileInfo
	size int64
}

func (i minifiedInfo) Size() int64 { return i.size }

// Compile-time proof that a minified tree still answers everything
// http.FileServer and fs.WalkDir ask of it.
var (
	_ fs.ReadDirFS = (*minifiedFS)(nil)
	_ fs.StatFS    = (*minifiedFS)(nil)
	_ io.Seeker    = (*minifiedFile)(nil)
)
