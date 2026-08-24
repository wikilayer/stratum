# Changelog

Stratum is a small CSS framework for the pages a Go server renders: tokens, layout, about twenty components and an icon sprite. It ships as a Go module holding one `fs.FS` you mount on a route and link stylesheets from; there is no build step and no CSS of your own. See the [README](README.md) for wiring it up and the [design-system reference](https://wikilayer.github.io/stratum/) for the markup each component takes.

Format: [Keep a Changelog](https://keepachangelog.com/). Versions are the tags a consumer pins with `go get github.com/wikilayer/stratum@vX.Y.Z`. Before 1.0 a minor bump may rename or remove a class; every such change is listed here with the markup to change.

## 0.2.0 — 2026-08-25

### Added

- **`body > footer`**, the closing strip: quiet type on the same edges as the bar at the top, a rule above it, links carrying the strip's colour. Put it last inside `<body>` and compose the row from `.cluster` as usual. Both full-window demos in the reference now carry one.

### Changed

- **`<body>` is a flex column of full window height, and `.content` takes the height left over.** This is what puts a footer on the floor of a page too short to reach it. Two things follow for pages built before this release, both visible rather than fatal:

  - A page whose `<body>` holds an element that relied on normal flow beside another — a floated block, two blocks meant to sit inline — now gets them stacked. Direct children of `<body>` are flex items.
  - `.content` no longer sets `min-height: calc(100vh - var(--header-h))`. Any page that leaned on that number for its own height should say so itself. In exchange, a page carrying a `.page-head-row` is no longer taller than the window by the height of that row, which had been showing up as a scrollbar on pages with nothing to scroll.

## 0.1.1 — 2026-08-24

### Fixed

- A `.block-id` longer than its 3em gutter wrapped, and at this line height the second line landed back on the first, so a five-digit id read as one smudge. It now runs left into the empty margin instead. No markup change.

## 0.1.0 — 2026-08-24

The first tagged release. There is no previous tag, so this entry is written for anyone who pinned a commit instead. It covers everything since [`df118a6`](https://github.com/wikilayer/stratum/commit/df118a6) (18 August 2026), and two renames from earlier in August that reach a pin older than that: `.page-tabs` became `.nav-tabs`, and `.caret` was dropped for the sprite's `chevron-down` icon at `.icon` size, which is what a disclosure arrow is drawn with now. For anything before that, read `github.com/wikilayer/stratum/compare/<your-pin>...v0.1.0`.

Most of these changes fail silently: nothing throws, the page just lays out differently. So work by search rather than by eye. `form-layout`, `split`, `tab-panel-`, `feed-action-`, `modal-body`, `--header-h`, `--aside-w`, `--leftnav-w`, `--content-max`, `column-wide`, `choice`, `field-label`, `row-value`, `is-revoked`, `danger-hover` between them find every site of the changes below, and `page-tabs`, `caret` the two above.

Do the `.split` step before the `.form-layout` rename, or the rename will hide which elements needed it.

### Changed

- **`.form-layout` is gone; `.split` is the two-column page.** They were one grid written twice: `.form-layout` put its fixed column on the right and made it 200px, `.split` put its own on the left and made it 280px. Both are now `.split`, fixed column on the right at `14rem`, and:

  - **First**, add `.split-side-first` to every `.split` you already have, or its side column moves to the other edge.
  - **Then** rename: `form-layout` → `split`, `form-layout-main` → `split-main`, `form-layout-aside` → `split-side`. These keep their side on the right and need no modifier.
  - The side column is `14rem` (224px) wherever it came from. To keep the width you had, set `--side-w` on the `.split` element itself: `12.5rem` for the old form layout, `17.5rem` for the old split.

  `.split-main` caps no line length of its own, where `.form-layout-main` capped at 32em. If you are coming from the form layout, add `.measure`, the utility that caps at the same 32em — unchanged, only now through the `--measure-form` token. `.measure-prose` is its counterpart for a column of running text:

  ```html
  <form class="split" style="--side-w: 12.5rem">
    <div class="split-main measure">…</div>
    <div class="split-side">…</div>
  </form>
  ```

- **Tabs pair by position instead of by id.** The component used to require the ids `#tab-1`…`#tab-4` and the classes `.tab-panel-1`…`-4`, which meant two tab groups on one page could not both work. Now the Nth radio pairs with the Nth label and the Nth panel, and the ids are yours to name.

  A panel is now `<section class="tab-panel">`, in place of the numbered classes. It has to be a `<section>` and a direct child of `.tabs`, because the rule counts sections among the siblings, and that is what keeps the `<div class="tabs-bar">` out of the count. A panel left as a `<div>` will not show at all. Keep the order of the example below — radios, then the bar, then the panels — and keep to four tabs: the rule is written out four times, and a fifth tab's panel never shows.

  Before:

  ```html
  <div class="tabs">
    <input type="radio" id="tab-1" name="mine" checked>
    <input type="radio" id="tab-2" name="mine">
    <div class="tabs-bar" role="tablist">
      <label for="tab-1" role="tab">First</label>
      <label for="tab-2" role="tab">Second</label>
    </div>
    <div class="tab-panel-1" role="tabpanel">…</div>
    <div class="tab-panel-2" role="tabpanel">…</div>
  </div>
  ```

  After — the radios and labels keep their `id`/`for` pairing, under any name; only the panels change:

  ```html
  <div class="tabs">
    <input type="radio" id="login-oauth" name="mine" checked>
    <input type="radio" id="login-password" name="mine">
    <div class="tabs-bar" role="tablist">
      <label for="login-oauth" role="tab">First</label>
      <label for="login-password" role="tab">Second</label>
    </div>
    <section class="tab-panel" role="tabpanel">…</section>
    <section class="tab-panel" role="tabpanel">…</section>
  </div>
  ```

- **Feed action variants are named for what an entry did:** `.feed-action-add`, `.feed-action-edit`, `.feed-action-remove`, replacing `.feed-action-INSERT` / `-UPDATE` / `-DELETE`. If your labels come from SQL statement names, map them to the three in your own code.

- **`.modal-body` no longer spaces its children apart.** Every dialog whose body holds more than one element needs `class="modal-body stack"`, or the elements sit flush against each other. `.stack` is the utility that puts a gap between siblings.

- **The page shell is sized in rem, not px.** `--header-h` (was 56px, now 3.5rem), `--aside-w` (256px → 16rem), `--leftnav-w` (220px → 13.75rem), `--content-max` (1100px → 68.75rem). The avatar, the switch and the map embed moved the same way, in their own rules rather than in tokens, so there is nothing to convert there.

  Root font is 16px unless the reader has changed it, so each of these still computes to the pixel value it had and a page that overrides none of them looks the same. **If you override any of these four, convert your own value to rem as well.** A px override still parses and still wins, which means the release fixes the defaults and leaves your page with the defect: the bar holds its height while the text inside it grows with the reader's font setting, and the contents spill out of it.

- `.column-wide` is `--measure-prose` (42em) rather than 720px, so the widest reading column grows with the reader's font too. `em` here is the body size, 17px, so this is 714px against the previous 720px.

- Three classes gave way to ones that already did their job:

  - `.field-label` → `.row-label`, the label class the README's form examples use. `.field-label` was documented nowhere.
  - `.is-revoked` → `.muted` on the row: the same dimming, under a name not borrowed from one application's domain.
  - `.row-value` → the read-only list the framework already had, which is a `<dl>` of label/value pairs rather than a single class on a value:

    ```html
    <dl class="row-inline-list">
      <div class="row-inline"><dt>Email:</dt><dd>alex@example.com</dd></div>
      <div class="row-inline"><dt>Joined:</dt><dd>3 March 2026</dd></div>
    </dl>
    ```

### Added

- `--measure-form` (32em) beside the existing `--measure-prose` (42em). These are the two line lengths content is capped at: a column of form controls, and running text.
- `.split-side-first` puts the fixed column on the leading edge (the left, in a left-to-right script). `--side-w` sets its width.
- A tab being moved to by keyboard is now visible: the radio that holds the focus is positioned off-screen, so the label draws the focus ring on its behalf.

### Fixed

- Dark mode reached through the operating system's setting was missing its shadow values, so dropdowns, dialogs and the switch drew shadows meant for a white page and had no visible edge on a dark one. Only a reader who had chosen dark by hand got the right ones.
- `.block-id`, the id shown in the gutter beside an addressable section, was dimmed to a 2.0:1 contrast ratio against the page, under the 4.5:1 that readable text needs. It is no longer dimmed.
- A heading that opens an `<article>` no longer carries a top margin. On a page whose text all sits in child sections, the wrapper around them is empty, so the first heading's margin collapsed through it and landed on top of the container's own padding: 56px of blank space where a page that opens with a paragraph has 24px and words.

### Removed

Nothing here has a replacement class; the renames are listed under Changed above.

- `--danger-hover` — a colour no rule in the framework read. If one of yours did, declare it yourself.
- `.choice-group`, `.choice`, `.choice-name`, `.choice-desc` — a fieldset of radio cards, reached by nothing. The same shape composes out of what remains:

  ```html
  <fieldset class="stack">
    <legend>Visibility</legend>
    <label class="card cluster">
      <input type="radio" name="visibility" value="public" checked>
      <span>
        <strong>Public</strong>
        <span class="fine-print">Anyone with the link can read it.</span>
      </span>
    </label>
  </fieldset>
  ```
