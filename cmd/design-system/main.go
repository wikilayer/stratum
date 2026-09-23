// Command design-system serves and builds Stratum's reference pages.
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/wikilayer/stratum"
)

const (
	designDir    = "design-system"
	sharedLayout = "_document.gohtml"
)

type renderer struct {
	source    fs.FS
	assetBase string
}

func (r renderer) render(name string) (resource, error) {
	if path.Ext(name) != ".html" || !fs.ValidPath(name) {
		return resource{}, fs.ErrNotExist
	}
	info, err := fs.Stat(r.source, name)
	if err != nil {
		return resource{}, err
	}
	tmpl, err := template.New(name).Funcs(template.FuncMap{
		"asset": func(asset string) string {
			return r.assetBase + strings.TrimPrefix(asset, "/")
		},
	}).ParseFS(r.source, sharedLayout, name)
	if err != nil {
		return resource{}, err
	}
	var body bytes.Buffer
	if err := tmpl.ExecuteTemplate(&body, name, nil); err != nil {
		return resource{}, err
	}
	return resource{
		name:         name,
		contentType:  "text/html; charset=utf-8",
		cacheControl: "no-cache",
		modTime:      info.ModTime(),
		body:         body.Bytes(),
	}, nil
}

type resource struct {
	name         string
	contentType  string
	cacheControl string
	modTime      time.Time
	body         []byte
}

func readResource(source fs.FS, name string) (resource, error) {
	if !fs.ValidPath(name) {
		return resource{}, fs.ErrNotExist
	}
	body, err := fs.ReadFile(source, name)
	if err != nil {
		return resource{}, err
	}
	info, err := fs.Stat(source, name)
	if err != nil {
		return resource{}, err
	}
	contentType := mime.TypeByExtension(path.Ext(name))
	if contentType == "" {
		contentType = http.DetectContentType(body)
	}
	return resource{
		name:         path.Base(name),
		contentType:  contentType,
		cacheControl: "public, max-age=0, must-revalidate",
		modTime:      info.ModTime(),
		body:         body,
	}, nil
}

func newHandler(source fs.FS) http.Handler {
	pages := renderer{source: source, assetBase: "/static/"}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/design-system/", http.StatusTemporaryRedirect)
	})
	mux.HandleFunc("GET /design-system/{name...}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if name == "" {
			name = "index.html"
		}
		var item resource
		var err error
		switch path.Ext(name) {
		case ".html":
			item, err = pages.render(name)
		case ".gohtml":
			err = fs.ErrNotExist
		default:
			item, err = readResource(source, name)
		}
		writeResource(w, r, item, err)
	})
	mux.HandleFunc("GET /static/{name...}", func(w http.ResponseWriter, r *http.Request) {
		item, err := readResource(stratum.Static, r.PathValue("name"))
		writeResource(w, r, item, err)
	})
	return mux
}

func writeResource(w http.ResponseWriter, r *http.Request, item resource, err error) {
	if errors.Is(err, fs.ErrNotExist) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sum := sha256.Sum256(item.body)
	etag := `"` + hex.EncodeToString(sum[:]) + `"`
	w.Header().Set("Cache-Control", item.cacheControl)
	w.Header().Set("Content-Type", item.contentType)
	w.Header().Set("ETag", etag)
	http.ServeContent(w, r, item.name, item.modTime, bytes.NewReader(item.body))
}

func build(source fs.FS, destination string) error {
	destination = filepath.Clean(destination)
	if destination == "." || destination == string(filepath.Separator) {
		return fmt.Errorf("refusing unsafe output directory %q", destination)
	}
	if err := os.RemoveAll(destination); err != nil {
		return err
	}
	pages := renderer{source: source, assetBase: "../static/"}
	if err := fs.WalkDir(source, ".", func(name string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() {
			return walkErr
		}
		if path.Ext(name) == ".gohtml" {
			return nil
		}
		var item resource
		var err error
		if path.Ext(name) == ".html" {
			item, err = pages.render(name)
		} else {
			item, err = readResource(source, name)
		}
		if err != nil {
			return err
		}
		return writeFile(filepath.Join(destination, "design-system", filepath.FromSlash(name)), item.body)
	}); err != nil {
		return err
	}
	if err := copyFS(stratum.Static, filepath.Join(destination, "static")); err != nil {
		return err
	}
	return writeFile(filepath.Join(destination, "index.html"), []byte(`<!doctype html>
<meta charset="utf-8">
<title>stratum — design system</title>
<meta http-equiv="refresh" content="0; url=design-system/">
<link rel="canonical" href="design-system/">
`))
}

func copyFS(source fs.FS, destination string) error {
	return fs.WalkDir(source, ".", func(name string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() {
			return walkErr
		}
		body, err := fs.ReadFile(source, name)
		if err != nil {
			return err
		}
		return writeFile(filepath.Join(destination, filepath.FromSlash(name)), body)
	})
}

func writeFile(name string, body []byte) error {
	if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		return err
	}
	return os.WriteFile(name, body, 0o644)
}

func openBrowser(url string) error {
	var command *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		command = exec.Command("open", url)
	case "windows":
		command = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		command = exec.Command("xdg-open", url)
	}
	return command.Start()
}

func run() error {
	addr := flag.String("addr", "127.0.0.1:8765", "HTTP listen address")
	buildSite := flag.Bool("build", false, "render the static site into dist and exit")
	open := flag.Bool("open", false, "open the design-system index in the default browser")
	flag.Parse()

	source := os.DirFS(designDir)
	if *buildSite {
		return build(source, "dist")
	}
	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		return err
	}
	url := "http://" + listener.Addr().String() + "/design-system/"
	log.Printf("design system: %s", url)
	if *open {
		if err := openBrowser(url); err != nil {
			return err
		}
	}
	return http.Serve(listener, newHandler(source))
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
