// stratum/cmd/icons regenerates static/icons.svg from the manifest at
// static/icons.txt. The design system and consumers use the same sprite.
//
// Manifest line forms:
//   <name>            — fetched from Lucide via unpkg.
//   <name> | <url>    — fetched from any URL (used for brand icons
//                        Lucide doesn't ship, e.g. github).
//
// Each source SVG is parsed: viewBox is preserved, and the styling
// (outline-stroke vs filled) is auto-detected from the presence of
// `stroke=` on the root element. Lucide icons all have stroke attrs
// so they render as outlines; Simple Icons brand SVGs have neither
// stroke nor fill at the root and render as solid currentColor.
//
// Run via `make icons-sync`. Re-runnable safely.

package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
)

const (
	manifestPath = "static/icons.txt"
	spritePath   = "static/icons.svg"
	lucideBase   = "https://unpkg.com/lucide-static/icons/"
)

type spec struct {
	name string
	url  string
}

func main() {
	specs, err := readManifest(manifestPath)
	if err != nil {
		fail("read manifest: %v", err)
	}
	if len(specs) == 0 {
		fail("manifest %s is empty", manifestPath)
	}

	symbols := make([]symbol, 0, len(specs))
	for _, s := range specs {
		body, err := fetch(s.url)
		if err != nil {
			fail("fetch %s (%s): %v", s.name, s.url, err)
		}
		sym, err := buildSymbol(s.name, body)
		if err != nil {
			fail("build %s: %v", s.name, err)
		}
		symbols = append(symbols, sym)
	}

	var sprite bytes.Buffer
	if err := writeSprite(&sprite, symbols); err != nil {
		fail("render sprite: %v", err)
	}
	if err := os.WriteFile(spritePath, sprite.Bytes(), 0o644); err != nil {
		fail("write external sprite: %v", err)
	}
	fmt.Printf("wrote %s — %d icons\n", spritePath, len(symbols))
}

//go:embed sprite.svg.tmpl
var spriteSource string

var spriteTmpl = template.Must(template.New("sprite").Parse(spriteSource))

type symbol struct {
	ID      string
	ViewBox string
	Style   string
	Inner   string
}

func writeSprite(w io.Writer, symbols []symbol) error {
	return spriteTmpl.Execute(w, symbols)
}

func readManifest(path string) ([]spec, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []spec
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		out = append(out, parseLine(line))
	}
	return out, scanner.Err()
}

func parseLine(line string) spec {
	parts := strings.SplitN(line, "|", 2)
	name := strings.TrimSpace(parts[0])
	if len(parts) == 1 {
		return spec{name: name, url: lucideBase + name + ".svg"}
	}
	return spec{name: name, url: strings.TrimSpace(parts[1])}
}

func fetch(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

var (
	rootSVGRe   = regexp.MustCompile(`(?s)<svg([^>]*)>(.*)</svg>`)
	viewBoxRe   = regexp.MustCompile(`viewBox="([^"]+)"`)
	strokeAttRe = regexp.MustCompile(`\bstroke=`)
	whitespace  = regexp.MustCompile(`\s+`)
	titleTagRe  = regexp.MustCompile(`(?s)<title>.*?</title>`)
)

const (
	outlineAttrs = `fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"`
	filledAttrs  = `fill="currentColor"`
)

// buildSymbol turns a fetched SVG into a single <symbol> line.
// Preserves the source viewBox; picks outline-stroke or filled
// styling based on whether the source SVG has a stroke attribute
// on its root element.
func buildSymbol(name, body string) (symbol, error) {
	m := rootSVGRe.FindStringSubmatch(body)
	if m == nil {
		return symbol{}, fmt.Errorf("no <svg>...</svg>")
	}
	rootAttrs := m[1]
	inner := m[2]

	vb := viewBoxRe.FindStringSubmatch(rootAttrs)
	if vb == nil {
		return symbol{}, fmt.Errorf("no viewBox")
	}
	viewBox := vb[1]

	style := filledAttrs
	if strokeAttRe.MatchString(rootAttrs) {
		style = outlineAttrs
	}

	// Strip <title> — accessible name comes from aria-label on the
	// host <svg> element, not from the symbol.
	inner = titleTagRe.ReplaceAllString(inner, "")
	inner = whitespace.ReplaceAllString(strings.TrimSpace(inner), " ")

	return symbol{ID: name, ViewBox: viewBox, Style: style, Inner: inner}, nil
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
