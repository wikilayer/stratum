package stratum

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRailIdentity_DemoLinksMarkAndNameTogether(t *testing.T) {
	demo, err := os.ReadFile("design-system/layout-3-rail.html")
	require.NoError(t, err)
	html := string(demo)
	start := strings.Index(html, `class="rail-identity"`)
	require.NotEqual(t, -1, start, "the three-rail demo has no rail identity")
	end := strings.Index(html[start:], "</a>")
	require.NotEqual(t, -1, end, "the rail identity is not a link")
	link := html[start : start+end]
	require.Contains(t, link, `class="rail-identity-mark`, "the home link must contain the mark")
	require.Contains(t, link, `class="rail-identity-name"`, "the home link must contain the name")
}

func TestRailIdentity_ReservesAFullWidthSquareForAnImageOrInitials(t *testing.T) {
	css, err := os.ReadFile("static/css/components/rail.css")
	require.NoError(t, err)
	styles := string(css)
	for _, want := range []string{
		".rail-identity-mark",
		"width: 100%",
		"height: auto",
		"aspect-ratio: 1",
		"img.rail-identity-mark",
		"object-fit: cover",
	} {
		if !strings.Contains(styles, want) {
			t.Errorf("rail identity styles do not contain %q", want)
		}
	}
}

func TestRailCurrentItem_UsesTheSameLeftMarkerAsContents(t *testing.T) {
	css, err := os.ReadFile("static/css/components/rail.css")
	require.NoError(t, err)
	styles := string(css)
	require.Contains(t, styles, ".rail-section li a[aria-current]")
	require.GreaterOrEqual(t, strings.Count(styles, "border-left-color: var(--accent)"), 2,
		"a current rail item and the current contents item must both draw the accent left marker")
}

func TestThreeRailShell_KeepsTheLeadingColumnsClose(t *testing.T) {
	tokens, err := os.ReadFile("static/css/base/tokens.css")
	require.NoError(t, err)
	require.Contains(t, string(tokens), "--rail-gutter: var(--space-5)",
		"the distance around page rails must be one adjustable token")
	css, err := os.ReadFile("static/css/base/layout.css")
	require.NoError(t, err)
	styles := string(css)
	for _, want := range []string{
		".content-with-head:has(> nav.leftnav) > nav.leftnav",
		"padding-inline-end: var(--rail-gutter)",
		".content-with-head:has(> nav.leftnav) > :where(.page-head-row, main)",
		"padding-inline-start: 0",
		"padding-inline-start: var(--rail-gutter)",
		".content-with-head > :is(main, aside)",
		"padding-block-start: var(--rail-gutter)",
	} {
		if !strings.Contains(styles, want) {
			t.Errorf("three-rail spacing does not contain %q", want)
		}
	}
}
