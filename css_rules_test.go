package stratum

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// What the framework promises about its own stylesheets, checked
// against them rather than left to whoever edits one next. A consumer
// overrides a component from its own sheet, and every one of these
// keeps that possible.

// Specificity stays flat, or a utility stops being able to override a
// component and the consumer's own rule stops being able to override
// either. An id selector and an !important each end that.
func TestCSS_KeepsSpecificityFlat(t *testing.T) {
	idSelector := regexp.MustCompile(`(^|[\s,>+~])#[a-zA-Z]`)
	forEachStylesheet(t, func(t *testing.T, path, css string) {
		for i, line := range strings.Split(css, "\n") {
			code, _, _ := strings.Cut(line, "/*")
			if strings.Contains(code, "!important") {
				t.Errorf("%s:%d: !important -- a consumer cannot override this from their own sheet", path, i+1)
			}
			selectorBeforeOpeningBrace, _, isRule := strings.Cut(code, "{")
			if isRule && idSelector.MatchString(selectorBeforeOpeningBrace) {
				t.Errorf("%s:%d: id selector in %q -- outranks every class a consumer could write", path, i+1, strings.TrimSpace(selectorBeforeOpeningBrace))
			}
		}
	})
}

// A rule outside a layer beats every rule inside one, whatever the
// order, so a stylesheet that forgets its @layer silently wins against
// the whole framework and against the consumer's own layers.
func TestCSS_EveryStylesheetDeclaresItsLayer(t *testing.T) {
	layers := strings.Split(strings.TrimSuffix(strings.TrimPrefix(cssLayerOrder, "@layer "), ";"), ", ")
	forEachStylesheet(t, func(t *testing.T, path, css string) {
		if path == "stratum.css" {
			return
		}
		for _, layer := range layers {
			if strings.Contains(css, "@layer "+layer+" {") {
				return
			}
		}
		t.Errorf("%s puts its rules in no layer, so they outrank every layered rule in the framework", path)
	})
}

// A var() naming a property nothing declares resolves to nothing, and
// the rule holding it does nothing -- silently, which is how a token
// removed from tokens.css leaves a component unstyled.
func TestCSS_EveryCustomPropertyUsedIsDeclared(t *testing.T) {
	declared := map[string]bool{}
	used := map[string]string{}
	declares := regexp.MustCompile(`(--[a-z0-9-]+)\s*:`)
	uses := regexp.MustCompile(`var\((--[a-z0-9-]+)`)

	forEachStylesheet(t, func(t *testing.T, path, css string) {
		for _, m := range declares.FindAllStringSubmatch(css, -1) {
			declared[m[1]] = true
		}
		for _, m := range uses.FindAllStringSubmatch(css, -1) {
			used[m[1]] = path
		}
	})

	for name, path := range used {
		if !declared[name] {
			t.Errorf("%s reads %s, which no stylesheet declares", path, name)
		}
	}
}

// The framework describes itself and never the applications built on
// it: a class or a comment naming one is a copy of that application's
// vocabulary, kept here where it goes stale unseen.
func TestCSS_NamesNoApplication(t *testing.T) {
	forbidden := []string{"wikilayer", "/settings", "/manage", "/owner", "/login"}
	forEachStylesheet(t, func(t *testing.T, path, css string) {
		lower := strings.ToLower(css)
		for _, word := range forbidden {
			if strings.Contains(lower, word) {
				t.Errorf("%s names %q -- the framework knows nothing of the pages built with it", path, word)
			}
		}
	})
}

func TestCSS_TreeNavigationNestsInEitherRail(t *testing.T) {
	css := readCSS(t, "css/components/rail.css")
	if !strings.Contains(css, ":where(.content > aside, .content > nav.leftnav) ul") {
		t.Fatal("rail list reset must stay less specific than nested navigation in either rail")
	}
	if !strings.Contains(css, ".toc ul ul,\n    .tree-nav ul ul") {
		t.Fatal("tree navigation and the table of contents must share their nested-list shape")
	}
}

func TestCSS_PageHeadKeepsItsTabsOnTheSharedBaseline(t *testing.T) {
	css := readCSS(t, "css/components/page-head.css")
	if !strings.Contains(css, ".page-head .nav-tabs") || !strings.Contains(css, "align-items: flex-end") {
		t.Fatal("page-head tabs must keep compact controls against the title's baseline")
	}
}

func TestCSS_CurrentTabAndOpenMenuShareSelectedSurface(t *testing.T) {
	css := readCSS(t, "css/components/nav-tabs.css")
	if !strings.Contains(css, ".nav-tabs a[aria-current],") ||
		!strings.Contains(css, ".nav-tabs .menu-host[open] > .menu-toggle") ||
		!strings.Contains(css, "background: var(--bg-hover)") {
		t.Fatal("current tabs and open menus must share the selected surface")
	}
}

func forEachStylesheet(t *testing.T, check func(t *testing.T, path, css string)) {
	t.Helper()
	err := fs.WalkDir(os.DirFS("static"), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".css") {
			return err
		}
		b, err := os.ReadFile(filepath.Join("static", path))
		if err != nil {
			return err
		}
		t.Run(path, func(t *testing.T) { check(t, path, string(b)) })
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
