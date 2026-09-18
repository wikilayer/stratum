package stratum

import (
	"os"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func designSystemPage(t *testing.T) *html.Node {
	t.Helper()
	f, err := os.Open("design-system/index.html")
	if err != nil {
		t.Fatalf("open the design-system page: %v", err)
	}
	defer f.Close()
	page, err := html.Parse(f)
	if err != nil {
		t.Fatalf("parse the design-system page: %v", err)
	}
	return page
}

func attrOf(n *html.Node, name string) string {
	for _, a := range n.Attr {
		if a.Key == name {
			return a.Val
		}
	}
	return ""
}

func hasClass(n *html.Node, want string) bool {
	for _, class := range strings.Fields(attrOf(n, "class")) {
		if class == want {
			return true
		}
	}
	return false
}

func forEachMenuItem(t *testing.T, visit func(item *html.Node, choice bool)) {
	t.Helper()
	var walk func(n *html.Node, inMenu, choice bool)
	walk = func(n *html.Node, inMenu, choice bool) {
		if n.Type == html.ElementNode {
			if attrOf(n, "role") == "menu" {
				inMenu, choice = true, hasClass(n, "dropdown-choice")
			}
			if inMenu && hasClass(n, "item") {
				visit(n, choice)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c, inMenu, choice)
		}
	}
	walk(designSystemPage(t), false, false)
}

func TestMenu_EveryItemSaysItIsOne(t *testing.T) {
	seen := 0
	forEachMenuItem(t, func(item *html.Node, _ bool) {
		seen++
		switch attrOf(item, "role") {
		case "menuitem", "menuitemradio", "menuitemcheckbox":
		default:
			t.Errorf("an item in a menu carries role %q: a menu whose children are not menuitems "+
				"is announced as an empty one, and the item as a stray button",
				attrOf(item, "role"))
		}
	})
	if seen == 0 {
		t.Fatal("no menu items found, so this test read nothing")
	}
}

func TestMenu_ACurrentChoiceIsCheckedRatherThanCurrent(t *testing.T) {
	radios := 0
	forEachMenuItem(t, func(item *html.Node, choice bool) {
		if current := attrOf(item, "aria-current"); current != "" {
			t.Errorf("an item in a menu marks itself aria-current=%q: that is how a link says "+
				"which view you are on, and a reader in a menu listens for aria-checked", current)
		}
		if !choice || attrOf(item, "role") != "menuitemradio" {
			return
		}
		radios++
		switch attrOf(item, "aria-checked") {
		case "true", "false":
		default:
			t.Errorf("a menuitemradio carries aria-checked=%q: without it a reader is told the "+
				"item is selectable and never told whether it is selected",
				attrOf(item, "aria-checked"))
		}
	})
	if radios == 0 {
		t.Fatal("no choice items found, so this test read nothing")
	}
}
