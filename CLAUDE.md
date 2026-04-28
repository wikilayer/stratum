# stratum — notes for Claude

Standalone embeddable UI framework (CSS + JS + icons). Lives at `~/Developer/stratum` as its own Go module. **Does not** depend on wikilayer (or any other consumer) — and must not.

- Module path: `github.com/wikilayer/stratum`.
- Public surface: a single `stratum.Static fs.FS` exported from `stratum.go`. Nothing else.
- Live preview: `https://wikilayer.github.io/stratum/` — rebuilt from `design-system/` on every push to `main` (see `.github/workflows/pages.yml`).

## Prime directive

**No page-specific classes.** If a name only fits one page of one app (`.login`, `.profile-grid`, `.consent-actions`), it's a missed abstraction. Promote it to a generic primitive (`.column-narrow`, `.split`, `.card-grid`, `.feed-*`) and let the template compose pages from those.

When you reach for a new class name, ask: "what would this be called if the same shape showed up on a different page tomorrow?" If the honest answer is a different name, pick that instead.

## Principles

1. **CSS comments stay short.** Prose (why a component exists, how to use it, markup conventions) lives in `README.md`. Inline CSS comments are reserved for one-line WHYs that aren't obvious from the rule itself.
2. **Every colour / spacing / radius / size goes through `var(--…)`** from `static/css/base/tokens.css`. A hex literal in a component rule is a bug — extend the tokens.
3. **`@layer` cascade order** is declared in `static/style.css`: `reset, tokens, base, layout, components, utilities`. Wrap every rule in its layer. Don't reorder.
4. **No preprocessors.** Native CSS Custom Properties + `@layer` + a sprinkle of `color-mix()`. Browser baseline: anything that ships `@layer` (Chrome 99 / Firefox 97 / Safari 15.4).
5. **One `<link rel="stylesheet">` in the consumer template** → `style.css`. It `@import`s the rest; HTTP/2 multiplexing handles the fan-out.

## Layout

```
stratum/
├── stratum.go              ← exports Static fs.FS — the only Go file
├── go.mod
├── Makefile                ← icons-sync, design-system targets
├── README.md               ← user-facing docs
├── CLAUDE.md               ← this file
├── .github/workflows/      ← Pages deploy
├── static/
│   ├── style.css           ← @layer order + @imports, no prose
│   ├── icons.svg / icons.txt / icons.LICENSE.txt
│   ├── theme.js            ← theme switcher (cookie + data-theme)
│   ├── copy.js             ← copy-to-clipboard for .url-pill
│   └── css/
│       ├── base/           ← reset, tokens, typography, layout
│       ├── components/     ← one file per component
│       └── utilities.css
├── design-system/
│   └── index.html          ← standalone showcase, opens via file://
└── cmd/
    └── icons/              ← icon-sprite generator (Lucide + Simple Icons)
```

## Adding a component

1. New file `static/css/components/<name>.css`. Wrap rules in `@layer components { … }`. Use a generic, shape-or-role-based name.
2. `@import` it from `static/style.css` in the components block.
3. Document the markup convention in `README.md` (Components section).
4. Add a live example in `design-system/index.html`.

Before introducing a new class, **read README.md first** — there is often a primitive you can compose with a modifier (`.button-sm`, `.alert-error`, `.column-narrow`) instead of a new class.

## Adding a utility / token

- Utility (single-purpose helper like `.text-center`, `.muted`) → `static/css/utilities.css`, `@layer utilities`.
- New colour / spacing / radius → `static/css/base/tokens.css`, `:root` plus the matching `html[data-theme="dark"]` override.

## Icons

Manifest: `static/icons.txt`. One `name` per line (Lucide), or `name | url` to fetch from any URL (used for brand icons Lucide doesn't ship — Simple Icons CC0). After editing the manifest run `make icons-sync` (delegates to `cmd/icons`); the sprite is rewritten at `static/icons.svg`. Licence text in `static/icons.LICENSE.txt`.

In templates: `<svg class="icon"><use href="/static/icons.svg#<name>"/></svg>`.

## Don't

- Don't depend on any consumer's package (no `wikilayer/...` imports). stratum is its own product.
- Don't introduce CSS classes named after a single consumer page (`.login`, `.consent-actions`, `.profile-grid`, etc.). See the prime directive.
- Don't use `!important`, `id` selectors, or deep nesting. Keep specificity flat so utilities reliably override components.
- Don't hand-edit `static/icons.svg` — it's regenerated from the manifest.
- Don't add a build step, npm package, or CSS preprocessor. Native everything.
- Don't write Russian (or any non-English) in committed files. Public repo, English only.

## How consumers wire it up

```go
import "github.com/wikilayer/stratum"

http.Handle("/static/*", http.StripPrefix("/static/",
    http.FileServer(http.FS(stratum.Static))))
```

Templates link `/static/style.css`, `/static/icons.svg#<name>`, `/static/theme.js`, `/static/copy.js`. If the consumer also serves its own files at `/static/`, compose two `fs.FS` together (mergedFS pattern — local files win on collision, missing names fall through to stratum).

## GitHub Pages

`design-system/index.html` is the human-readable spec. The workflow at `.github/workflows/pages.yml` builds a `dist/` containing both `design-system/` and `static/` (so the relative `../static/style.css` link inside `index.html` resolves), and publishes via `actions/deploy-pages`. The `dist/index.html` at the root is a tiny redirect to `design-system/`.

To enable the first time:
```bash
gh api repos/wikilayer/stratum/pages -X POST -f build_type=workflow
```

## Tagging a release

When the API is stable enough that a consumer can drop the `replace` directive:

```bash
git tag v0.1.0
git push --tags
```

Consumers then `go get github.com/wikilayer/stratum@v0.1.0` and remove their local `replace`.
