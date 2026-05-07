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

// CSSAssets lists every CSS file in the bundle in cascade order. Hosts
// can either link static/style.css (one request that chains N more
// @import requests, all render-blocking and serial) or iterate this
// slice to emit N parallel <link rel="stylesheet"> tags — the latter
// is dramatically faster on cold loads with HTTP/2 multiplexing.
//
// When iterating, also inline CSSLayerOrder in a <style> block before
// the first <link> so layer precedence is fixed regardless of which
// stylesheet finishes loading first.
var CSSAssets = []string{
	"css/base/tokens.css",
	"css/base/reset.css",
	"css/base/typography.css",
	"css/base/layout.css",
	"css/components/button.css",
	"css/components/field.css",
	"css/components/avatar.css",
	"css/components/url-pill.css",
	"css/components/callout.css",
	"css/components/map.css",
	"css/components/card.css",
	"css/components/table.css",
	"css/components/breadcrumb.css",
	"css/components/dropdown.css",
	"css/components/aside.css",
	"css/components/article.css",
	"css/components/page-head.css",
	"css/components/tabs.css",
	"css/components/alert.css",
	"css/components/modal.css",
	"css/components/feed.css",
	"css/utilities.css",
}

// CSSLayerOrder is the @layer declaration that needs to ship before
// any of the CSSAssets so the cascade resolves deterministically. Drop
// it into a tiny inline <style> in <head> to avoid taking another
// network round-trip just to declare layer precedence.
const CSSLayerOrder = "@layer reset, tokens, base, layout, components, utilities;"

func mustSub(f fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(f, dir)
	if err != nil {
		panic(err)
	}
	return sub
}
