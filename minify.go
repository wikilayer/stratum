package stratum

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"path"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/js"
)

type minifiedFS struct {
	source    fs.FS
	processed map[string][]byte
}

func newMinifiedFS(source fs.FS) fs.FS {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/javascript", js.Minify)

	out := &minifiedFS{source: source, processed: map[string][]byte{}}
	err := fs.WalkDir(source, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		mediaType := ""
		switch path.Ext(name) {
		case ".css":
			mediaType = "text/css"
		case ".js":
			mediaType = "text/javascript"
		default:
			return nil
		}
		raw, err := fs.ReadFile(source, name)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		if err := m.Minify(mediaType, &buf, bytes.NewReader(raw)); err != nil {
			return err
		}
		out.processed[name] = buf.Bytes()
		return nil
	})
	if err != nil {
		panic("stratum: minify assets: " + err.Error())
	}
	styles, err := joinAssets(cssLayerOrder, cssAssets, "\n", out.processed)
	if err != nil {
		panic("stratum: bundle stratum.css: " + err.Error())
	}
	scripts, err := joinAssets("", jsAssets, ";\n", out.processed)
	if err != nil {
		panic("stratum: bundle stratum.js: " + err.Error())
	}
	out.processed["stratum.css"] = styles
	out.processed["stratum.js"] = scripts
	return out
}

func joinAssets(prefix string, names []string, separator string, assets map[string][]byte) ([]byte, error) {
	var out bytes.Buffer
	out.WriteString(prefix)
	for _, name := range names {
		body, ok := assets[name]
		if !ok {
			return nil, fmt.Errorf("%s is listed but not among the static files", name)
		}
		out.WriteString(separator)
		out.Write(body)
	}
	return out.Bytes(), nil
}

func (m *minifiedFS) Open(name string) (fs.File, error) {
	body, ok := m.processed[name]
	if !ok {
		return m.source.Open(name)
	}
	info, err := fs.Stat(m.source, name)
	if err != nil {
		sourceName := ""
		switch name {
		case "stratum.js":
			sourceName = jsAssets[0]
		default:
			return nil, err
		}
		info, err = fs.Stat(m.source, sourceName)
		if err != nil {
			return nil, err
		}
	}
	return &minifiedFile{
		Reader: bytes.NewReader(body),
		info: minifiedInfo{
			FileInfo: info,
			name:     path.Base(name),
			size:     int64(len(body)),
		},
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
	if _, ok := m.processed[name]; !ok {
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
	name string
	size int64
}

func (i minifiedInfo) Name() string { return i.name }
func (i minifiedInfo) Size() int64  { return i.size }

var (
	_ fs.ReadDirFS = (*minifiedFS)(nil)
	_ fs.StatFS    = (*minifiedFS)(nil)
	_ io.Seeker    = (*minifiedFile)(nil)
)
