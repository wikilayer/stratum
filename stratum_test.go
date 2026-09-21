package stratum

import (
	"io/fs"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCSSAssetsResolve checks every entry in cssAssets actually opens
// in the embedded FS. Catches drift if a component CSS gets renamed
// or removed without updating the manifest.
func TestCSSAssetsResolve(t *testing.T) {
	for _, name := range cssAssets {
		f, err := Static.Open(name)
		if err != nil {
			t.Errorf("cssAssets entry %q: %v", name, err)
			continue
		}
		_ = f.Close()
	}
}

// TestCSSAssetsCoverComponents walks css/components/ in the embedded FS
// and asserts every file is listed in cssAssets — so adding a new
// component CSS without registering it fails CI instead of silently
// shipping unused styles to consumers that follow the new pattern.
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
			t.Errorf("css file %q not in cssAssets — add it (in cascade order)", path)
		}
		return nil
	})
	require.NoError(t, err)
}

// Every framework helper belongs to the one-request bundle. A new script
// must not be silently served only as an individual, unregistered asset.
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
			t.Errorf("JavaScript file %q not in jsAssets", path)
		}
		return nil
	})
	require.NoError(t, err)
}
