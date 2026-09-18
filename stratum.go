// Package stratum is a minimal CSS framework designed to embed into
// Go projects. Tokens, a layered cascade, reusable components, an icon
// sprite, and a few tiny vanilla-JS helpers — vendored as a Go module
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
// Then link /static/stratum.css from your template. See README and the
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
// URLs from your templates: /static/stratum.css, /static/icons.svg, and
// /static/stratum.js (or an individual helper on a narrow page).
//
// Stylesheets and scripts are served minified: the sources carry the
// reasoning behind each rule and helper, which is worth its bytes to
// whoever edits them and nothing at all to a browser. stratum.css and
// stratum.js are true one-request bundles.
var Static = newMinifiedFS(mustSub(embedded, "static"))

var cssAssets = []string{
	"css/base/tokens.css",
	"css/base/reset.css",
	"css/base/typography.css",
	"css/base/layout.css",
	"css/components/button.css",
	"css/components/field.css",
	"css/components/avatar.css",
	"css/components/tag.css",
	"css/components/url-pill.css",
	"css/components/callout.css",
	"css/components/map.css",
	"css/components/card.css",
	"css/components/table.css",
	"css/components/scroll-x.css",
	"css/components/breadcrumb.css",
	"css/components/dropdown.css",
	"css/components/switch.css",
	"css/components/rail.css",
	"css/components/article.css",
	"css/components/page-head.css",
	"css/components/nav-tabs.css",
	"css/components/tabs.css",
	"css/components/alert.css",
	"css/components/modal.css",
	"css/components/feed.css",
	"css/components/entry.css",
	"css/utilities.css",
}

var jsAssets = []string{
	"theme.js",
	"copy.js",
	"modal.js",
	"dropdown.js",
	"rail.js",
	"autosubmit.js",
}

const cssLayerOrder = "@layer reset, tokens, base, layout, components, utilities;"

func mustSub(f fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(f, dir)
	if err != nil {
		panic(err)
	}
	return sub
}
