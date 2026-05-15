// Mobile collapse for table-of-contents <details> blocks.
//
// Pages with an inline TOC mark it up as
//   <details class="rail-section toc-section" open>
//     <summary class="rail-section-summary">…</summary>
//     <nav>…</nav>
//   </details>
//
// On desktop the rail is wide enough to show the TOC fully expanded,
// so we keep [open] forced on. On mobile (≤720px) the TOC sits
// inline above the article — having the full link list expanded
// pushes the actual content below the fold. We strip [open] so the
// summary acts as a tap target the reader can choose to expand.
//
// The toggle is layout-driven, not user-driven: we follow the media
// query both on initial load and on resize. User toggles within the
// current viewport are preserved (no event listener on the details
// element), only crossing the breakpoint resets the state.

(function () {
    var mq = window.matchMedia('(max-width: 720px)');

    function apply() {
        var nodes = document.querySelectorAll('details.toc-section');
        for (var i = 0; i < nodes.length; i++) {
            if (mq.matches) {
                nodes[i].removeAttribute('open');
            } else {
                nodes[i].setAttribute('open', '');
            }
        }
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', apply);
    } else {
        apply();
    }
    // matchMedia .addEventListener is the modern API; fall back to
    // .addListener for Safari < 14 if anyone's still on that.
    if (mq.addEventListener) {
        mq.addEventListener('change', apply);
    } else if (mq.addListener) {
        mq.addListener(apply);
    }
})();
