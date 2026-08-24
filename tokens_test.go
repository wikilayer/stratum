package stratum

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A theme has two ways in: the attribute a reader sets by hand, and
// the media query that follows their system. Both name the same
// palette and are written out separately, so the pair drifts silently
// -- the shadows did, and the reader who never touched the switch got
// the light ones on a dark background.
func TestDarkPaletteIsTheSameBothWaysIn(t *testing.T) {
	css := readCSS(t, "css/base/tokens.css")

	chosen := tokensIn(t, css, `html[data-theme="dark"]`)
	followed := tokensIn(t, css, `html:not([data-theme])`)

	for name, want := range chosen {
		got, ok := followed[name]
		if !ok {
			t.Errorf("%s is set for a reader who picks dark and left out for one whose system says so", name)
			continue
		}
		if got != want {
			t.Errorf("%s is %q when picked and %q when followed", name, want, got)
		}
	}
	for name := range followed {
		if _, ok := chosen[name]; !ok {
			t.Errorf("%s is set for a system in dark mode and left out for a reader who picks dark", name)
		}
	}
}

// readCSS reads a stylesheet from the source tree rather than through
// Static, which serves it minified: what is checked here is how the
// rule is written, and minifying drops the very spelling that is.
func readCSS(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("static", path))
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}

// tokensIn returns the custom properties the rule with this selector
// declares, by taking the text between its opening brace and the
// matching close. Nothing nests inside these blocks, so a brace count
// is enough to read them without a CSS parser.
func tokensIn(t *testing.T, css, selector string) map[string]string {
	t.Helper()
	at := strings.Index(css, selector)
	if at < 0 {
		t.Fatalf("no rule for %s", selector)
	}
	open := strings.Index(css[at:], "{")
	if open < 0 {
		t.Fatalf("%s opens no block", selector)
	}
	body := css[at+open+1:]
	depth := 1
	end := 0
	for i, r := range body {
		switch r {
		case '{':
			depth++
		case '}':
			depth--
		}
		if depth == 0 {
			end = i
			break
		}
	}
	if end == 0 {
		t.Fatalf("%s never closes", selector)
	}

	out := map[string]string{}
	for _, line := range strings.Split(body[:end], "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "--") {
			continue
		}
		name, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		value = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(value), ";"))
		if comment := strings.Index(value, "/*"); comment >= 0 {
			value = strings.TrimSpace(value[:comment])
		}
		out[strings.TrimSpace(name)] = value
	}
	return out
}
