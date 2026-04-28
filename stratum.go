// Package stratum is a minimal CSS framework designed to embed into
// Go projects. Tokens, a layered cascade, ~20 components, an icon
// sprite, and two tiny vanilla-JS helpers — vendored as a Go module
// so a server gets a styled UI by mounting one fs.FS and linking one
// stylesheet. No npm, no build step, no preprocessor.
//
// Wire it up:
//
//	import "github.com/wikilayer/stratum"
//
//	http.Handle("/static/*", http.StripPrefix("/static/",
//	    http.FileServer(http.FS(stratum.Static))))
//
// Then link /static/style.css from your template. See README and the
// design-system reference at https://wikilayer.github.io/stratum/.
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
