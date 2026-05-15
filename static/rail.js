// Per-section open/closed state for <details class="rail-section">.
// Each section names itself with data-rail-section="NAME".
// Cookie format: rail-state=name1:1,name2:0  (1=open, 0=closed).
// Server is expected to read this cookie and set the [open] attribute
// at render time so there's no flash; this file just keeps the cookie
// in sync when the user toggles a section.

(function () {
    var COOKIE = 'rail-state';

    function readState() {
        var match = document.cookie.match(new RegExp('(?:^|; )' + COOKIE + '=([^;]*)'));
        var state = {};
        if (!match) return state;
        match[1].split(',').forEach(function (pair) {
            var i = pair.indexOf(':');
            if (i > 0) state[pair.slice(0, i)] = pair.slice(i + 1) === '1';
        });
        return state;
    }

    function writeState(state) {
        var pairs = [];
        Object.keys(state).forEach(function (k) {
            pairs.push(k + ':' + (state[k] ? '1' : '0'));
        });
        var oneYear = 60 * 60 * 24 * 365;
        document.cookie = COOKIE + '=' + pairs.join(',') + ';path=/;max-age=' + oneYear + ';samesite=lax';
    }

    document.addEventListener('toggle', function (e) {
        var el = e.target;
        if (!(el instanceof HTMLElement) || !el.classList.contains('rail-section')) return;
        var name = el.getAttribute('data-rail-section');
        if (!name) return;
        var state = readState();
        state[name] = el.open;
        writeState(state);
    }, true);
})();
