// Click-outside auto-close for every <details class="menu-host">
// on the page. Native <details> only closes when the user clicks
// the same <summary> again — every other dropdown UI on the web
// closes when the user clicks elsewhere, and people expect that.
//
// One global pointerdown listener is enough: walk up from the
// click target, find each open .menu-host that isn't an ancestor
// of the click, close it. Cheap and stable across dynamically-
// inserted hosts because the listener is on document.

(function () {
    document.addEventListener('pointerdown', function (e) {
        document.querySelectorAll('details.menu-host[open]').forEach(function (host) {
            if (!host.contains(e.target)) {
                host.removeAttribute('open');
            }
        });
    });
})();
