package stratum

import (
	"io/fs"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBundle_RefusesAListedFileThatIsGone(t *testing.T) {
	_, err := joinAssets("", []string{"css/kept.css", "css/removed.css"}, "\n",
		map[string][]byte{"css/kept.css": []byte("a{}")})
	require.ErrorContains(t, err, "css/removed.css",
		"a stylesheet removed without the list following drops out of the bundle, and every page loses its rules unannounced")
}

func TestCSSAssetsCoverComponents(t *testing.T) {
	listed := make(map[string]bool, len(cssAssets))
	for _, n := range cssAssets {
		listed[n] = true
	}
	err := fs.WalkDir(Static, "css", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".css") {
			return nil
		}
		if !listed[path] {
			t.Errorf("css file %q not in cssAssets, so it ships to no consumer — add it (in cascade order)", path)
		}
		return nil
	})
	require.NoError(t, err)
}

func TestJSAssetsCoverScripts(t *testing.T) {
	listed := make(map[string]bool, len(jsAssets))
	for _, name := range jsAssets {
		listed[name] = true
		if _, err := fs.Stat(Static, name); err != nil {
			t.Errorf("jsAssets entry %q: %v", name, err)
		}
	}
	source := mustSub(embedded, "static")
	err := fs.WalkDir(source, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".js") {
			return err
		}
		if !listed[path] {
			t.Errorf("JavaScript file %q not in jsAssets, so it is left out of the one-request bundle", path)
		}
		return nil
	})
	require.NoError(t, err)
}
