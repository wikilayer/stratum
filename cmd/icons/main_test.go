package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSprite_DrawsOutlineAndFilledIconsAsSymbols(t *testing.T) {
	outline, err := buildSymbol("lock", `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" stroke="currentColor">
  <rect width="18" height="11" x="3" y="11"/>
</svg>`)
	require.NoError(t, err)
	filled, err := buildSymbol("github", `<svg role="img" viewBox="0 0 24 24"><title>GitHub</title><path d="M12 .3"/></svg>`)
	require.NoError(t, err)

	var sprite bytes.Buffer
	require.NoError(t, writeSprite(&sprite, []symbol{outline, filled}))

	require.True(t, strings.HasPrefix(sprite.String(), `<?xml version="1.0" encoding="UTF-8"?>`))
	require.Contains(t, sprite.String(), "\n"+`  <symbol id="lock" viewBox="0 0 24 24" `+outlineAttrs+`><rect width="18" height="11" x="3" y="11"/></symbol>`+"\n",
		"an icon whose root strokes is drawn as an outline in the colour around it")
	require.Contains(t, sprite.String(), "\n"+`  <symbol id="github" viewBox="0 0 24 24" `+filledAttrs+`><path d="M12 .3"/></symbol>`+"\n",
		"an icon with no stroke is filled, and its title is left to the page that shows it")
	require.True(t, strings.HasSuffix(sprite.String(), "</symbol>\n</svg>\n"))
}
