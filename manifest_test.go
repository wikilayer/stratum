package stratum

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// The source manifest and the in-memory bundle list the same sheets in
// the same order, so the design-system reference and consumers see the
// same cascade.
func TestStyleEntryAndAssetListAgree(t *testing.T) {
	entry := readCSS(t, "stratum.css")

	var imported []string
	for _, m := range regexp.MustCompile(`@import url\("([^"]+)"\)`).FindAllStringSubmatch(entry, -1) {
		imported = append(imported, m[1])
	}

	require.Len(t, imported, len(cssAssets), "stratum.css and cssAssets list different numbers of stylesheets")
	for i, path := range imported {
		if path != cssAssets[i] {
			t.Errorf("stylesheet %d: stratum.css has %q, cssAssets has %q -- cascade order differs between source and bundle",
				i, path, cssAssets[i])
		}
	}
}

// The layer order is declared twice for the same reason, and a host
// that inlines the constant while the file says something else gets a
// cascade that depends on which stylesheet finishes loading first.
func TestLayerOrderIsDeclaredTheSameBothWays(t *testing.T) {
	entry := readCSS(t, "stratum.css")
	first, _, _ := strings.Cut(entry, "\n")

	if strings.TrimSpace(first) != cssLayerOrder {
		t.Errorf("stratum.css opens with %q, cssLayerOrder is %q", strings.TrimSpace(first), cssLayerOrder)
	}
}
