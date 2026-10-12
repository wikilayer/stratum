package stratum

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStyleEntryAndAssetListAgree(t *testing.T) {
	entry := readCSS(t, "stratum.css")

	var imported []string
	for _, m := range regexp.MustCompile(`@import url\("([^"]+)"\)`).FindAllStringSubmatch(entry, -1) {
		imported = append(imported, m[1])
	}

	require.Len(t, imported, len(cssAssets),
		"stratum.css and cssAssets list different numbers of stylesheets, so the design-system reference and consumers see different cascades")
	for i, path := range imported {
		if path != cssAssets[i] {
			t.Errorf("stylesheet %d: stratum.css has %q, cssAssets has %q -- cascade order differs between source and bundle",
				i, path, cssAssets[i])
		}
	}
}

func TestLayerOrderIsDeclaredTheSameBothWays(t *testing.T) {
	entry := readCSS(t, "stratum.css")
	first, _, _ := strings.Cut(entry, "\n")

	if strings.TrimSpace(first) != cssLayerOrder {
		t.Errorf("stratum.css opens with %q, cssLayerOrder is %q -- a host inlining the constant gets a cascade that depends on which stylesheet loads first",
			strings.TrimSpace(first), cssLayerOrder)
	}
}
