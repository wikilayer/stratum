# stratum

### → [Live design-system reference](https://wikilayer.github.io/stratum/)

**A minimal CSS framework designed to embed into Go projects.**

Tokens, a layered cascade, ~20 components, an icon sprite, and two tiny vanilla-JS helpers — vendored as a single Go module. A Go server gets a styled UI by importing the package, mounting one `fs.FS`, and linking one stylesheet. No npm, no build step, no preprocessor.

Native CSS Custom Properties + `@layer` + a sprinkle of `color-mix()`. Works in any browser that ships `@layer` (Chrome 99 / Firefox 97 / Safari 15.4 — 2022+).

## Why it exists

Most CSS frameworks are either heavy (Bootstrap, Tailwind — with their own toolchain) or drop-in single-stylesheet kits (Pico, Simple.css) that stop short of components. Stratum sits in between: ~20 tiny components, plus tokens and utilities, sized for an internal app or a side-project where shipping a Node toolchain alongside a single Go binary is wrong. The design-system page is the spec — if you can't build a new page out of what's documented there, the answer is to add a primitive, not a page-specific class.

## Install (Go)

```go
import "github.com/wikilayer/stratum"
```

Mount the static FS under `/static/` (use it directly, or layer your own files on top via `fs.FS` composition):

```go
http.Handle("/static/*", http.StripPrefix("/static/",
    http.FileServer(http.FS(stratum.Static))))
```

Link from your base template:

```html
<link rel="stylesheet" href="/static/style.css">
<script src="/static/theme.js" defer></script>
<script src="/static/copy.js" defer></script>
```

That's it. Everything else is plain HTML + classes.

## Cascade order (`@layer`)

```
reset → tokens → base → layout → components → utilities
```

Declared in `static/style.css`. Every rule in the framework is wrapped in its layer. Layers later in the list win in the cascade — so utilities can always override a component, components can always override base typography, and so on.

## Tokens

All colours, spacing, radii, type sizes, shadows, durations live in `static/css/base/tokens.css` as CSS variables. There's a light default and a `html[data-theme="dark"]` override; `prefers-color-scheme: dark` falls through to the same dark values when no theme is explicitly chosen.

Hardcoding a hex inside a component is a smell. Add a token first.

Highlights:

| Variable                    | Meaning                                |
|-----------------------------|----------------------------------------|
| `--bg`, `--bg-elevated`, `--bg-subtle`, `--bg-hover`, `--bg-active` | Surface levels |
| `--fg`, `--fg-muted`, `--fg-subtle`, `--fg-on-accent` | Text levels |
| `--border`, `--border-muted` | Two contrast levels for separators |
| `--accent`, `--accent-hover`, `--accent-bg` | Brand accent colour set |
| `--danger*`, `--success-*`, `--callout-*` | Semantic colour roles |
| `--text-xs … --text-4xl` | Type scale (1.25 ratio, body 17px) |
| `--space-1 … --space-8` | 4-step spacing scale, in `rem` |
| `--radius-sm/md/lg/pill` | Corner radii |
| `--header-h`, `--aside-w`, `--content-max` | Page-shell sizes |
| `--duration-fast`, `--easing-out` | Motion |

## Base layer

Sits above tokens, below components. Three files:

- `reset.css` — `box-sizing: border-box`, zero default margins on text blocks, `font: inherit` on form controls, `display: block` on media. Nothing more.
- `typography.css` — body font, six heading levels with a 1.25-ratio scale, paragraphs with bottom margin, `<code>`/`<pre>`/`<blockquote>`. Headings are clean by default; opt into the Wikipedia / MDN baseline rule per element with `.h-rule`, or globally inside an `<article>`.
- `layout.css` — page chrome (sticky `body > header`, `.brand`, `.content` frame, optional `<aside>` rail). Plus reusable column primitives (`.column*`) and `.split` for two-column pages.

## Layout primitives

### `.column`, `.column-narrow`, `.column-wide`

Reading-width columns, centred. Pick one when a page wants to constrain its content to a single readable line.

```html
<section class="column-narrow">…</section>     <!-- 32em — auth, dialogs -->
<section class="column">…</section>             <!-- 38em — onboarding, prose -->
<article class="column-wide">…</article>        <!-- 720px — legal docs -->
```

### `.split`, `.split-side`, `.split-main`

Two-column page with a fixed-width side rail and a fluid main column. Below 720px collapses to one stack with the side moved on top.

```html
<div class="split">
  <aside class="split-side">…</aside>
  <section class="split-main">…</section>
</div>
```

## Components

Every component lives in its own file under `static/css/components/`. Markup conventions below.

### `.button`

```html
<a class="button button-primary"  href="…">…</a>
<a class="button button-secondary" href="…">…</a>
<a class="button button-ghost"     href="…">…</a>
<button class="button button-danger">…</button>
<a class="button button-primary button-sm" href="…">…</a>
```

Modifiers stack. Use `<a>` for navigation, `<button>` for actions. Icons go inline before the label.

### `.row`, `.field`, `.input`, `.choice-group`

Form primitives.

`.row` is a label-on-top stack. `.field` is the same but targets a single isolated field with a tighter label.

```html
<form class="stack">
  <div class="row">
    <label class="row-label" for="x">Display name</label>
    <input id="x" type="text" required>
  </div>
  <button class="button button-primary">Save</button>
</form>
```

`.row-inline-list` / `.row-inline` render a read-only `<dl>` of `Label: value` rows.

`.choice-group` / `.choice` / `.choice-name` / `.choice-desc` is a fieldset of mutually-exclusive radios, each option a clickable card with icon + name + description. Selected state uses `:has(input:checked)`.

`.form-layout` / `.form-layout-main` / `.form-layout-aside` is a generic two-column form (main + 200px sidebar slot) — used on `/settings` for the avatar.

Inputs / selects inside `.row` and `.field` get the framework's text-input look automatically. Custom `<select>` chevron is painted with two CSS gradients so it follows the theme.

### `.tabs`

Pure-CSS, radio-driven. Up to four tabs out of the box; extend the selector pairs in `tabs.css` for more.

```html
<div class="tabs">
  <input type="radio" id="tab-1" name="my-tabs" checked>
  <input type="radio" id="tab-2" name="my-tabs">
  <div class="tabs-bar" role="tablist">
    <label for="tab-1" role="tab">First</label>
    <label for="tab-2" role="tab">Second</label>
  </div>
  <div class="tab-panel-1" role="tabpanel">…</div>
  <div class="tab-panel-2" role="tabpanel">…</div>
</div>
```

### `.alert`

Inline message block. Variants: `.alert-error`, `.alert-success`. Always left-aligned (won't inherit `.text-center`). Use `.alert-success` as a one-shot banner after a redirect — same shape, no separate primitive needed.

### `.modal`

Dialog box on top of a backdrop. Built on the native `<dialog>` element — ESC closes, click on backdrop closes (wired by `modal.js`), focus trapped, no scroll-lock gymnastics. Three sub-zones: header / body / footer.

```html
<button class="button" data-modal-open="add-member">+ Add member</button>

<dialog id="add-member" class="modal">
  <header class="modal-header">
    <p class="modal-title">Add member</p>
    <button class="modal-close" data-modal-close aria-label="Close">×</button>
  </header>
  <div class="modal-body">…form fields…</div>
  <footer class="modal-footer">
    <button class="button" data-modal-close>Cancel</button>
    <button class="button button-primary">Add</button>
  </footer>
</dialog>
```

`data-modal-open="ID"` on any clickable opens the dialog with that id. `data-modal-close` on any clickable inside the dialog closes it. Load `/static/modal.js` once on the page.

### `.callout`

GitHub-flavoured callouts (`> [!NOTE]` rendered from markdown). Variants: `.callout-note`, `.callout-tip`, `.callout-important`, `.callout-warning`, `.callout-caution`. Wrap the title in `.callout-title` with an icon.

### `.map`, `.map-caption`

Embedded location iframe (rendered from `> [!MAP]` blockquotes in markdown). The container clips a single full-width `<iframe>` to a rounded card; an optional `.map-caption` strip sits below.

```html
<div class="map">
  <iframe src="https://maps.google.com/maps?q=51.4779,0.0015&z=15&output=embed" loading="lazy" referrerpolicy="no-referrer-when-downgrade" allowfullscreen></iframe>
  <div class="map-caption">Royal Observatory, Greenwich</div>
</div>
```

### `.card`, `.card-grid`, `.card-link`

Bordered surface for grouped content. `.card-grid` lays cards in a responsive grid. Add `.card-link` to an `<a>` that wraps the whole card — it gets the flex-row layout, link colour, and accent border on hover.

```html
<ul class="card-grid">
  <li><a class="card card-link" href="…"><svg class="icon">…</svg> Title</a></li>
</ul>
```

### `.avatar`, `.avatar-lg`

Round chip with initials, or `<img class="avatar">` for a Gravatar.

### `.url-pill`

Copy-this-code block — a `<code>` value paired with a compact `.copy-btn` (wired by `copy.js`). The name says URL because that's the canonical use, but the body is a generic `<code>` and works for any short string the user needs to paste somewhere.

```html
<div class="url-pill">
  <code>https://example.com/api</code>
  <button class="copy-btn" type="button"
          data-label-copy="Copy" data-label-copied="Copied">Copy</button>
</div>
```

### `.page-head`, `.page-tabs`

Title plus an icon-tab strip on the same baseline. The active tab's underline replaces the page-head's bottom border locally. Ideal for "this page has a Content / History / Settings switcher" layouts.

### `.breadcrumb`

```html
<nav class="breadcrumb">
  <a href="/">home</a><span class="sep">/</span>
  <a href="/x">Section</a><span class="sep">/</span>
  <span class="current">Current</span>
</nav>
```

### `.menu-host`, `.menu-toggle`, `.dropdown`

Click-to-reveal menu built on `<details>`. The header avatar dropdown is the canonical use; the panel inside hosts `.dropdown .item` rows and optional `.dropdown-section` blocks.

```html
<details class="menu-host">
  <summary class="menu-toggle">…</summary>
  <div class="dropdown" role="menu">
    <a class="item" href="…">Settings</a>
    <div class="dropdown-section">
      <div class="label">Theme</div>
      <div class="segmented">…</div>
    </div>
  </div>
</details>
```

### `.segmented`

Horizontal radio-style picker (button-group). Use `aria-pressed="true"` to mark the active option.

### `.feed`, `.feed-item`, `.feed-time`, `.feed-actor`, `.feed-action`, `.feed-target`

Flat list of timestamped activity entries. `.feed-action` carries a coloured tag — variants by suffix (`feed-action-INSERT`, `…-UPDATE`, `…-DELETE`; rename in your CSS if your domain uses different verbs). `.feed-target-gone` strikes through a target whose object no longer exists.

### `.data-table`, `.meta`

Generic admin / settings table. `.meta` is the same idea but for a `<dl>` of read-only `Label: value` facts.

### `.toc`, `.toc-link`, `.recent`

Right-rail widgets — table of contents with active-link highlight, recent-activity list with title + relative time.

### `.block`, `.block-id`

Addressable section inside long-form content. `<section class="block" id="…">` carries `scroll-margin-top` matching the sticky header; `.block-id` is the gutter `#`-link rendered by markdown post-processing.

## Utilities

Single-purpose helpers in `static/css/utilities.css`. Add new ones sparingly — most "I need this in two places" patterns belong in a component.

| Class            | Effect |
|------------------|--------|
| `.muted`         | secondary foreground colour |
| `.empty`         | italic + muted, for "no data" placeholders |
| `.fine-print`    | small + muted text, for disclaimers under forms |
| `.text-center`, `.text-left`, `.text-right` | `text-align` |
| `.stack > * + *` | vertical rhythm via margin-top (lobotomized owl) |
| `.cluster`       | horizontal flex with gap and wrap |
| `.inline-form`   | `display: inline` for inline POST forms |
| `.icon`          | 1em-square inline SVG, follows `currentColor` |
| `.h-rule`        | thin baseline rule under a heading (Wikipedia / MDN look). Headings are clean by default — opt in per `<h1>`/`<h2>`. Applies automatically to all `h1`/`h2` inside an `<article>` so rendered markdown gets the rule for free |
| `.eyebrow`       | small uppercase muted label — sidebar section titles, dropdown labels, and other "small caps over a list" places |

## Icons

The sprite at `static/icons.svg` is regenerated from `static/icons.txt`. Manifest format: one name per line (fetched from Lucide via unpkg), or `name | url` to fetch from any URL (used for brand icons Lucide doesn't ship — Simple Icons CC0).

```bash
make icons-sync
```

Rendering:

```html
<svg class="icon" aria-hidden="true">
  <use href="/static/icons.svg#globe"/>
</svg>
```

The sprite uses `stroke="currentColor"` for outline icons and `fill` for solid brand glyphs, so colour follows the surrounding text.

License: outline icons are Lucide / Feather (ISC + MIT subset). Brand icons are Simple Icons (CC0). Full text lives in `static/icons.LICENSE.txt` and ships with the sprite.

## JavaScript helpers

- `theme.js` — wires every `[data-theme-set]` button. Writes a `theme` cookie (`light` / `dark` / blank for auto), applies `data-theme` on `<html>` immediately, updates `aria-pressed` on the buttons. Server is responsible for reading the cookie on each request and emitting `<html data-theme="…">` on the initial render to avoid flash.
- `copy.js` — wires every `.copy-btn` inside a `.url-pill`. Copies the `<code>` value to clipboard, swaps the button label to `data-label-copied`, then back after 1500ms. Localised labels stay in templates, not in JS.
- `modal.js` — wires `[data-modal-open="ID"]` triggers and `[data-modal-close]` close-buttons. Uses native `<dialog>.showModal()` / `.close()`; adds click-on-backdrop-to-close on top of what the platform gives you for free.

All three are zero-dependency, ~30 lines each, safe to load with `defer`.

## Adding a component

1. New file in `static/css/components/<name>.css`. Wrap rules in `@layer components { … }`. Keep selectors flat — no `id` selectors, no deep nesting.
2. `@import` it from `static/style.css` in the components block.
3. Document the markup convention in this README under Components.
4. Add a live example to `design-system/index.html`.
5. Re-check imports use generic class names. If the name only fits one page of one app, you missed an abstraction — pick a shape-based or role-based name instead.

## Adding a token

`static/css/base/tokens.css`. Define under `:root` first, then add the dark-theme override in `html[data-theme="dark"]` and (when the value differs from light) the matching `prefers-color-scheme: dark` block. Reference via `var(--…)` everywhere else.

## Design system

`design-system/index.html` is a standalone reference — opens with `file://`, no server needed. It lives next to the framework so any change to a primitive can be sanity-checked alongside the docs in seconds.

```bash
make design-system   # opens it in the default browser
```

If you add a primitive and don't add an example here, future-you will reinvent it. Update the page.

## Constraints

Things this framework deliberately doesn't do:

- **No preprocessors.** Native CSS only. If you reach for Sass, the rule isn't generic enough.
- **No build step.** A single static folder, served as-is.
- **No JavaScript framework.** Two tiny `.js` files, both vanilla, both optional.
- **No `!important`, no `id` selectors, no deep nesting.** Specificity stays flat so utilities reliably override components.
- **No page-specific classes.** If a name only fits one page (`.login`, `.profile-grid`, `.consent-actions`), it's the wrong abstraction. Compose pages from the primitives above.

## Layout

```
stratum/
├── stratum.go              ← the only Go file: exports Static fs.FS
├── go.mod
├── Makefile                ← icons-sync, design-system targets
├── README.md
├── CLAUDE.md               ← notes for assistants working on this package
├── static/
│   ├── style.css           ← entry: @layer order + @imports
│   ├── icons.{svg,txt,LICENSE.txt}
│   ├── theme.js
│   ├── copy.js
│   ├── modal.js
│   └── css/
│       ├── base/{tokens,reset,typography,layout}.css
│       ├── components/*.css
│       └── utilities.css
├── design-system/
│   └── index.html          ← live reference, file://-friendly
└── cmd/
    └── icons/              ← icon-sprite generator (Lucide + Simple Icons)
```
