package stratum

import (
	"regexp"
	"strings"
	"testing"
)

// A stylesheet is registered in two places: style.css, which a host
// may link on its own, and CSSAssets, which a host iterates to link
// them in parallel. A component reaching only one of them is styled
// for half the hosts, and the half that misses out is the
// design-system reference, which is where the framework is read.
func TestStyleEntryAndAssetListAgree(t *testing.T) {
	entry := readCSS(t, "style.css")

	var imported []string
	for _, m := range regexp.MustCompile(`@import url\("([^"]+)"\)`).FindAllStringSubmatch(entry, -1) {
		imported = append(imported, m[1])
	}

	if len(imported) != len(CSSAssets) {
		t.Fatalf("style.css imports %d stylesheets, CSSAssets lists %d", len(imported), len(CSSAssets))
	}
	for i, path := range imported {
		if path != CSSAssets[i] {
			t.Errorf("stylesheet %d: style.css has %q, CSSAssets has %q -- cascade order differs between the two ways in",
				i, path, CSSAssets[i])
		}
	}
}

// The layer order is declared twice for the same reason, and a host
// that inlines the constant while the file says something else gets a
// cascade that depends on which stylesheet finishes loading first.
func TestLayerOrderIsDeclaredTheSameBothWays(t *testing.T) {
	entry := readCSS(t, "style.css")
	first, _, _ := strings.Cut(entry, "\n")

	if strings.TrimSpace(first) != CSSLayerOrder {
		t.Errorf("style.css opens with %q, CSSLayerOrder is %q", strings.TrimSpace(first), CSSLayerOrder)
	}
}
