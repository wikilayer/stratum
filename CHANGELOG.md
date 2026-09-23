# Changelog

Stratum is a small CSS framework for the pages a Go server renders: tokens, layout, reusable components and an icon sprite. It ships as a Go module holding one `fs.FS` you mount on a route and link stylesheets from; there is no build step, and a consuming page is meant to reach its whole look through these class names rather than write a stylesheet of its own. See the [README](README.md) for wiring it up and the [design-system reference](https://wikilayer.github.io/stratum/) for the markup each component takes.

Format: [Keep a Changelog](https://keepachangelog.com/). Versions are the tags a consumer pins with `go get github.com/wikilayer/stratum@vX.Y.Z`. Before 1.0 any release may rename or remove a class, or need a change to your markup, the patch digit included; a release that does opens with **Not drop-in** and says what to change.

## 0.4.21 — 2026-09-23

Drop-in. Existing markup keeps working.

### Added

- `.tree-nav` renders nested, animated navigation with compact indentation, one disclosure chevron per expandable row and the standard accent bar on the active link.
- `.sequence-nav` keeps previous and next document links in equal halves and clamps long titles to two lines; when one neighbour is absent, its half remains empty.
- The design-system reference includes complete three-rail and modal examples.

### Fixed

- Three-rail documents stop the reading column from collapsing at narrow widths, and the mobile navigation control stays under the pointer when its drawer opens.
- Modal fields follow the active colour theme; modal headers and footers use the same height and the same padding on both sides.

## 0.4.20 — 2026-09-23

Drop-in. Existing markup keeps working.

### Fixed

- Cards inside the same `.card-grid` row stretch to the row height instead of ending at different baselines when their content lengths differ.

## 0.4.19 — 2026-09-23

**Not drop-in.** Replace `.wiki-card-mark` with the generic media-card structure: put `.media-card` on the linked card, `.media-card-visual` on its 4:3 image or initials, and wrap the inline icon and title in `.media-card-body`.

### Added

- `.media-card`, `.media-card-visual` and `.media-card-body` form a vertical linked card whose artwork or initials sit above an inline status icon and wrapping title.
- The design-system card section shows initials, artwork, long-title and alternate-state media cards.

### Changed

- Avatar colour slot 2 uses the shared yellow palette in both themes and keeps dark ink for contrast.

### Removed

- `.wiki-card-mark`, which described one consumer instead of a reusable component.

## 0.4.18 — 2026-09-23

**Not drop-in.** Remove `.block-id` elements: Stratum no longer provides a visible gutter permalink component. Keep the address on the `.block` itself if fragments still need to resolve.

### Added

- `.wiki-card-mark` turns the existing image-or-initials element into a compact 4:3 thumbnail for a horizontal card; no consumer CSS is required.

### Changed

- `.rail-identity-mark` uses the same 4:3 artwork proportion as wiki cards instead of forcing a square. Images continue to fill the frame with `object-fit: cover`, so their edges may be cropped; existing markup keeps working.

### Removed

- `.block-id` and its layout-specific exceptions. Whether an addressable section exposes a permalink is a product decision, not something Stratum infers from the surrounding page shell.

## 0.4.17 — 2026-09-22

Drop-in. Existing markup keeps working.

### Fixed

- A first block rendered with `.editable-heading` starts at the same height as its ordinary reading heading.

## 0.4.16 — 2026-09-22

Drop-in. Existing markup keeps working.

### Fixed

- Hidden children no longer leave a stack gap behind, including error messages that only appear after a failed form submission.

## 0.4.15 — 2026-09-22

Drop-in. Existing markup keeps working.

### Changed

- Modal headers and footers share the same compact vertical rhythm, while modal bodies keep equal padding on all sides.
- Fields fill a modal body, and a stack owns the spacing between its fields instead of adding the field margin twice.
- Textareas inside `.row` or `.field` receive the same theme-aware treatment as other controls.

### Added

- The design system includes a standalone page exercising short, form and long-form modals.

## 0.4.14 — 2026-09-22

Drop-in. Existing markup keeps working.

### Added

- `.editable-heading` makes a heading and its pencil mark one accessible edit target, with the mark sized to the lowercase letters at every heading level.
- `.modal-compact` and `.modal-body-compact` provide tighter dialog chrome and body spacing.
- The icon sprite includes the Lucide `pencil` icon.

### Fixed

- Page controls beside a right rail take the width they need and grow into the title column instead of wrapping onto a second row.

## 0.4.10 — 2026-09-21

Drop-in. Existing markup keeps working. The optional `.heading-text` wrapper
limits an article H2 accent rule to its title when tags follow it.

### Added

- Article H2 headings use a thicker accent rule that ends with the title instead of crossing the reading column.
- `.heading-text` lets that rule stop before a trailing `.tag-row`.

### Fixed

- Long page titles stay inside the reading column instead of colliding with the language and page controls.
- Long left navigation columns grow with the document instead of becoming clipped nested scrollers.

## 0.4.9 — 2026-09-21

Drop-in. No markup or consumer CSS change is required.

### Fixed

- A `.rail-identity-mark` that also carries `.avatar` now resets the avatar's fixed height, so an existing image or initials renderer fills the full-width square without consumer CSS.

## 0.4.8 — 2026-09-21

Drop-in. Existing markup keeps working. Visually verify pages where a
`.content-with-head` directly contains `nav.leftnav`, `main` and `aside`,
because their gutters are now shared and slightly tighter.

### Added

- `.rail-identity`, `.rail-identity-mark` and `.rail-identity-name` form a linked identity above a navigation rail. The square mark accepts either an image or coloured initials, and the name below it leads home through the same link.
- `--rail-gutter` controls the space between document columns and below their shared page head. Its default is `1.5rem`.

### Fixed

- A current navigation-rail link marked with `aria-current` now carries the same accent bar as the active table-of-contents link.
- In the three-rail shell, the page title and article now start at the left edge of their horizontal rule, while the article and right rail share the same top gutter.

## 0.4.7 — 2026-09-20

### Fixed

- On a document shorter than the viewport, `.content-with-head` now lets its navigation rails fill the remaining page height, so their vertical rules stay joined to the title rule and reach the footer.

## 0.4.6 — 2026-09-20

### Fixed

- `.content-with-head` keeps its title row and document columns content-sized, so the horizontal and vertical rules meet when leading navigation is absent.
- The two-column design-system example now exercises the same title controls as the three-column shell.

## 0.4.5 — 2026-09-20

### Changed

- `.content-with-head` now drops an absent `nav.leftnav` from its grid instead of reserving an empty leading column. The same document shell therefore covers pages with and without leading navigation.
- The three-column page shell grows to `76rem` before it stops, giving the reading column its full `--measure-prose` width while keeping both navigation columns fixed.

## 0.4.2 — 2026-09-18

No runtime changes. The public documentation now points at the current release and Go Reference, and the repository checks documentation comments as part of its complete build.

## 0.4.1 — 2026-09-13

### Added

- `.link-button`: an action that reads as a line of text — muted, underlined, no box — for the rare one that sits beside prose or a field of metadata and would shout as a button. It stays a `<button>`, because an action dressed as an `<a href>` is one a reader can open in a new tab and never reach.

  ```html
  <dd>25 April 2026 · <button type="button" class="link-button">Delete account</button></dd>
  ```

## 0.4.0 — 2026-09-12

**Not drop-in.** Two edits to your markup, and they are independent of each other.

**Every `.dropdown-choice` panel**, whether or not it says `role="menu"`: the check-mark is drawn from `aria-checked="true"` where it used to be drawn from `aria-current="true"`. The selector is `.dropdown-choice .item[aria-checked="true"]`, so the attribute alone brings the mark back; the roles below have nothing to do with it. Swap the attribute and delete the old one — a leftover `aria-current` in a menu is read out as a claim about which view you are on.

**Every panel that says `role="menu"`**: give each `.item` a role. This changes nothing you can see, and everything a screen reader hears.

### Changed

- A panel that says `role="menu"` promises a shape a screen reader reads by, and its items were plain buttons carrying no role: the menu was announced as holding nothing, and each item as a stray control. Give an item that picks one of a set `role="menuitemradio"`, an item that does something `role="menuitem"`, and write `aria-checked="false"` on the choices that are not current — without it a reader is told an item is selectable and never told it is not selected. A `.dropdown-separator` is not an item and takes `role="separator"`.

  Before, and after:

  ```html
  <button class="item" aria-current="true">Editor</button>
  <button class="item">Viewer</button>
  ```

  ```html
  <button class="item" role="menuitemradio" aria-checked="true">Editor</button>
  <button class="item" role="menuitemradio" aria-checked="false">Viewer</button>
  ```

  A whole panel, choices and an acting item either side of a separator:

  ```html
  <div class="dropdown dropdown-choice" role="menu">
    <button class="item" role="menuitemradio" aria-checked="true">Editor</button>
    <button class="item" role="menuitemradio" aria-checked="false">Viewer</button>
    <div class="dropdown-separator" role="separator"></div>
    <button class="item" role="menuitem">
      <svg class="icon" aria-hidden="true"><use href="/static/icons.svg#trash-2"/></svg>
      Remove
    </button>
  </div>
  ```

## 0.3.4 — 2026-09-11

### Changed

- In a `.dropdown-choice` panel the check-mark, any leading `.icon`, and every item's label sit `--space-4` (1rem) from the panel edge instead of `--space-3` (0.75rem), so the mark stands as far from that edge as the label does from the other one. This touches every choice panel, not only the ones using the leading icon 0.3.3 added, and nothing needs editing either way. For the leading icon it is the release that makes it sit right: 0.3.3 alone puts it visibly too close to the edge.

## 0.3.3 — 2026-09-11

**Not drop-in.** The row-list fix below needs a change to your markup; without it those links lose their hover underline. The leading icon this release adds is misplaced until 0.3.4, so take the two together.

### Added

- `.dropdown-choice` keeps a column on every item for the check-mark that marks the current choice. An item may now put an `.icon` in that column instead, so an item that does something — remove, hand over, sign out — begins at the same edge as the items that choose.
- `.dropdown-separator`, new in this version, is a one-pixel rule you place between two items. It is for setting an acting item apart from the choices above it; `.dropdown-section`, which was the only grouping there was, indents what it holds and so would push that item in from the edge the others span.

  ```html
  <div class="dropdown dropdown-choice" role="menu">
    <button class="item" aria-current="true">Editor</button>
    <button class="item">Viewer</button>
    <div class="dropdown-separator" role="separator"></div>
    <button class="item">
      <svg class="icon" aria-hidden="true"><use href="/static/icons.svg#trash-2"/></svg>
      Remove
    </button>
  </div>
  ```

- `trash-2` and `crown` join the icon sprite: a bin for removing somebody, a crown for making somebody the owner.

### Fixed

- `.row-list-item-link:hover` drew its underline across the whole link, avatar included, so the rule crossed the circle and struck the initial inside it. The underline is now drawn on the link's child elements, leaving out `.avatar` and `.icon`. **Every such link needs its text inside an element** — a bare text node is not a child and gets no underline at all, avatar or no avatar:

  ```html
  <a class="row-list-item-link" href="…">
    <span class="avatar">A</span>
    <span>A Reader</span>
  </a>
  ```

## 0.3.2 — 2026-09-04

### Changed

- Where a `.content` follows a `.page-head-row`, the gap between the two grows from `--space-5` (1.5rem) to `--space-6` (2rem), so the head and what hangs below it read as one page rather than two. Affects only pages using `.page-head-row`; everything else is unchanged, and nothing needs editing.

## 0.3.1 — 2026-09-01

### Removed

- `toc.js` is gone; **delete any `<script>` tag pointing at it**, or every page render asks for a file that is no longer served. It changed a `<details>` state after first paint according to viewport width, moving the article and creating cumulative layout shift.

  Whether a section starts open is now your server's decision, taken at render time: set the `open` attribute yourself. `rail.js` keeps the reader's choice in a cookie named `rail-state`, holding `name:1` or `name:0` per section, comma-separated, where the name is each `<details class="rail-section">`'s own `data-rail-section`. Read that cookie and render `open` from it. If you relied on toc.js to open the contents on a wide screen, that behaviour is gone and there is no client-side replacement: decide it on the server or leave the sections closed.

## 0.3.0 — 2026-09-01

**Not drop-in.** Link `stratum.css` and `stratum.js`, and stop using the Go surfaces removed below, or the build stops compiling.

### Added

- `stratum.css` is the whole framework in one stylesheet, and `stratum.js` the whole of its browser helpers in one script: link those two and a page makes one request for each instead of a chain of `@import`. The separate minified helpers are still served for a page that wants only one behaviour.

### Removed

- `style.css`, and the exported `CSSAssets` and `CSSLayerOrder`. Anything in your Go code naming either stops compiling; link `stratum.css` instead. A pre-1.0 cleanup, taken while WikiLayer is the only consumer.

## 0.2.5 — 2026-08-28

### Added

- `.cluster-center`, beside `.cluster-spread` and `.cluster-baseline`: the row centred in what holds it. For a row that is the whole of its container rather than one side of it — a strip of links closing a page, a pair of store badges — where against the left edge the first item runs under the rounded corner of a phone screen.

## 0.2.4 — 2026-08-28

### Added

- The Apple mark in the icon sprite, beside the GitHub and Google ones, for a row of sign-in buttons that names all three.

## 0.2.3 — 2026-08-26

### Fixed

- The bar at the top and the strip at the bottom sat on 1.5rem of side padding on a phone while `<main>` between them sat on 2rem, so the brand, the breadcrumb trail and the footer links all began a few pixels left of every line of text on the page. Both now take the column's own gutter as their floor; above the width where the centring calc takes over nothing changes. No markup change.

## 0.2.2 — 2026-08-26

### Fixed

- On a phone the page rail sat on 1rem of side padding while `<main>` below it sat on 2rem, so the page had two left edges and the rail's headings started further out than the text they introduce. The rail now takes the column's padding, and the line dividing it from the content is drawn inside that padding instead of as a border on the box: a border ran the full width of the screen, while every other rule on the page — under the tabs, under a heading — stops at the column's edge. No markup change.

## 0.2.1 — 2026-08-25

### Fixed

- `.page-head-row` collapsed to the width of its own text in 0.2.0, taking the page title and its tabs to the middle of the window. It centred itself with `margin: 0 auto` and let the page decide its width, which works in normal flow and does not in the flex column 0.2.0 gave `<body>`: auto margins on a flex item absorb the free space, and an item that named no width of its own shrank to its content. It now sets `width: 100%` and the band is a row across the page again. No markup change; if you copied the pattern onto a body child of your own, it needs the same line.

## 0.2.0 — 2026-08-25

The changes below are silent: no page fails to render, it lays out differently.

All of it turns on the shell this framework assumes, so here it is stated once: `<body>` holds the site's chrome and one `.content` element, `.content` being the centred frame that holds `<main>` and any rail beside it. `body > header` is the top bar, `body > footer` from this release the closing strip, and `.content` the column between them. What you work on is the template that writes your `<body>` tag and everything it puts directly inside it.

Before upgrading, screenshot each of your templates. One of the effects below shows up on nothing but a before-and-after comparison, and taking the "before" once is cheaper than re-pinning 0.1.1 to get it back.

None of this adds a stylesheet or a mount point. It all ships in `css/base/layout.css`, which `style.css` already imports and `CSSAssets` already lists, `CSSAssets` being the slice of stylesheet paths a host iterates to link them itself.

### Added

- **`body > footer`, the closing strip.** If your pages already end with a `<footer>` there, read the Changed entry below first: this claims it. A `<footer>` that is a direct child of `<body>` renders as site chrome: 15px type (`--text-sm`) in the muted foreground (`--fg-muted`), 24px above and below (`--space-5`), a 1px `--border-muted` rule on top, and links inheriting the strip's colour instead of the link colour, going to full-strength `--fg` on hover and keeping the framework's focus ring. Its horizontal padding is the one `body > header` uses, `max(var(--space-5), calc((100% - var(--content-max)) / 2 + var(--space-6)))`, where `--space-6` is 2rem: 24px on a narrow window, and from a little before `--content-max` onward an inset that aligns the strip with the text inside the content column rather than with the column's own edge, 2rem further out. The `100%` is the body's content width, not the window's, so a scrollbar shifts it by its own width. Override `--content-max` above both of them, on `:root` or a shared ancestor, and the bar and the strip move together; set on the footer it moves the strip alone. `--space-5` is the vertical padding too, so changing it there moves the strip's left edge off the bar's on a narrow window.

  ```html
  <body>
    <header>…</header>
    <div class="content">…</div>
    <footer>
      <div class="cluster">
        <span>An example site</span>
        <a href="/source">Source</a>
        <a href="/contact">Contact</a>
      </div>
    </footer>
  </body>
  ```

  The two full-window shell demos in the [reference](https://wikilayer.github.io/stratum/) each carry that markup, rendered.

### Changed

- **A `<footer>` already sitting directly inside `<body>` picks up that strip**, wanted or not: its type shrinks to 15px, its colour dims, a 1px rule appears above it. Check the ones your pages carry. The selector is `body > footer` and nothing else, so a footer inside a card, a modal or `<main>` keeps its own shape, and wrapping yours in a `<div>` is how a page opts out of the strip entirely.

  Redefining `--text-sm`, `--fg-muted`, `--border-muted` or `--space-5` on the footer element re-sizes, re-colours or re-spaces the strip, `--border-muted: transparent` included, which leaves the border in place, invisible, still holding its 1px of height. They are the ordinary tokens, not footer-scoped ones, so every descendant that reads them follows: a `.fine-print` line inside a footer whose `--fg-muted` you changed changes with it.

- **`<body>` is a flex column of at least window height** (`min-height: 100vh` and then `min-height: 100dvh`, the second winning wherever it parses, so a page longer than the window still scrolls the way it did and a phone's retracting address bar does not add a screenful), **and `.content` takes the height left over** (`flex: 1`, which is `1 1 0%`: in a column that basis supersedes a `height` you set on `.content` yourself, while a `min-height` of yours still clamps it from below, and its automatic minimum keeps it from shrinking below its content until you say `min-height: 0` or give it an `overflow` other than `visible`). That is what rests the footer on the floor of a page too short to reach it. `.content` has to stay a direct child of `<body>` for that: it is the flex item being stretched.

  Every in-flow direct child of `<body>` on every page is a flex item now, whether or not that page has a `.content` or a footer. A child with `position: absolute` or `fixed` — an overlay, a toast, the usual absolutely-positioned skip link — is out of flow and so not a flex item, and none of the six points below reach it. A `display: contents` wrapper is the opposite case and worth checking: it generates no box of its own, so its children become the flex items, one per child.

  There is nothing to grep for here, but the list is one line in the console of any page:

  ```js
  [...document.body.children].map(n => [n.tagName, n.className, getComputedStyle(n).position])
  ```

  For every row that is not `absolute` or `fixed`, six things follow:

  - Two children that used to sit side by side, a pair meant to be inline, stack instead. Wrap that pair, and only that pair, in a `<div>` of your own: inside it they are back in normal flow, and the wrapper takes their place in the column. Do not wrap `.content` along with them.
  - A child that used to be only as wide as its contents — a bare `<table>`, `<img>` or `<button>` sitting straight in `<body>` — now spans the full width, because a flex column stretches its items across. Give it `align-self: start` to keep the old width; on an `<img>` that also restores the height its aspect ratio had been growing along with the width.
  - `float` stops applying to such a child at all, because a flex item is never floated. Same wrapper, same reason.
  - Margins between body children no longer collapse, into each other or through the body's own edge, so two stacked blocks that each carried a positive vertical margin now sit the sum of both apart rather than the larger of the two. This is the one that hides: nothing looks broken, the page has simply loosened, which is why the screenshots above are worth taking. The fix is to drop the margin on one of the two.
  - An inline element left in the flow straight inside `<body>` is blockified as a flex item, so it takes a line of its own. Wrap it to keep it inline with what it sat next to.
  - Bare text sitting straight in `<body>`, outside any element, becomes an anonymous flex item on its own line. Wrap it in a `<p>` if its position mattered. Whitespace alone does not, so indentation in your template costs nothing.

  Watch the auto-margin case in particular: a body child that centred itself with `margin-inline: auto` and took its width from the page now shrinks to its content, because auto margins on a flex item absorb the free space instead of leaving it to the item. Give such a child `width: 100%` and it claims the width first, centring what is left over. A child with a width of its own keeps it and stays centred as before. The column adds no `gap` of its own, so the only new space between children is the uncollapsed margin above, and `body > header` keeps its height and its sticky behaviour as a flex item.

  `.content` no longer carries `min-height: calc(100vh - var(--header-h))`, because taking the leftover height of a full-window column already does what that number did, and that line alone needs no action. A `.content` nested inside something else — a `<main>`, a wrapper `<div>` — is not the flex item and does not stretch. Hoist it to be a direct child of `<body>`: giving the wrapper `flex: 1` stretches the wrapper and leaves `.content` sized to its content inside it. On a page built without `.content`, the six points above apply the same way, and nothing stretches: the footer follows the content instead of resting on the floor of the window. Give whichever block should absorb the slack `flex: 1` to get that behaviour back.

### Fixed

- A page carrying a `.page-head-row`, the band holding a page's title and its row of `.nav-tabs` above `.content`, was taller than the window by the height of that band, so a page with nothing to scroll still scrolled by that much. What fixes it is `.content` losing its own `min-height: calc(100vh - var(--header-h))`, a subtraction that counted the top bar and not the band; the column above sizes the page instead, and a page whose content is shorter than the window is exactly as tall as the window. The reset this framework ships keeps `<body>` free of margin, which is what the last part depends on: a margin of your own there is outside the `100dvh` and brings the overflow back. Padding and a border are not, since the reset also sets `box-sizing: border-box`, and they sit inside it.

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
