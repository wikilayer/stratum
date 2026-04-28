package stratum

import (
	"embed"
	"io/fs"
)

//go:embed static
var embedded embed.FS

// Static is the embedded asset tree rooted at static/. Mount it under
// /static/ (or wherever) on your HTTP server, and link the produced
// URLs from your templates: /static/style.css, /static/icons.svg,
// /static/theme.js, /static/copy.js.
var Static = mustSub(embedded, "static")

func mustSub(f fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(f, dir)
	if err != nil {
		panic(err)
	}
	return sub
}
