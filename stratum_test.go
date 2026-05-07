package stratum

import (
	"io/fs"
	"strings"
	"testing"
)

// TestCSSAssetsResolve checks every entry in CSSAssets actually opens
// in the embedded FS. Catches drift if a component CSS gets renamed
// or removed without updating the manifest.
func TestCSSAssetsResolve(t *testing.T) {
	for _, name := range CSSAssets {
		f, err := Static.Open(name)
		if err != nil {
			t.Errorf("CSSAssets entry %q: %v", name, err)
			continue
		}
		_ = f.Close()
	}
}

// TestCSSAssetsCoverComponents walks css/components/ in the embedded FS
// and asserts every file is listed in CSSAssets — so adding a new
// component CSS without registering it fails CI instead of silently
// shipping unused styles to consumers that follow the new pattern.
func TestCSSAssetsCoverComponents(t *testing.T) {
	listed := make(map[string]bool, len(CSSAssets))
	for _, n := range CSSAssets {
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
			t.Errorf("css file %q not in CSSAssets — add it (in cascade order)", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
