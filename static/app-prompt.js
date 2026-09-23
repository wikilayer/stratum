(function () {
    var storagePrefix = 'stratum:app-prompt:';

    function mobilePlatform() {
        var agent = navigator.userAgent || '';
        var iPad = navigator.platform === 'MacIntel' && navigator.maxTouchPoints > 1;
        if (/iPhone|iPad|iPod/.test(agent) || iPad) return 'ios';
        if (/Android/.test(agent)) return 'android';
        return '';
    }

    function isDismissed(key) {
        try {
            return Number(localStorage.getItem(storagePrefix + key)) > Date.now();
        } catch (_) {
            return false;
        }
    }

    function dismiss(prompt) {
        var days = Number(prompt.dataset.dismissDays) || 7;
        var expires = Date.now() + days * 24 * 60 * 60 * 1000;
        try {
            localStorage.setItem(storagePrefix + (prompt.dataset.dismissKey || 'default'), String(expires));
        } catch (_) {}
        prompt.hidden = true;
    }

    function forwardLocation(link) {
        if (!link.hasAttribute('data-current-location')) return;
        var target = new URL(link.href);
        target.pathname = location.pathname;
        target.search = location.search;
        target.hash = location.hash;
        link.href = target.href;
    }

    document.querySelectorAll('[data-app-prompt]').forEach(function (prompt) {
        var platform = prompt.dataset.platform || mobilePlatform();
        var platforms = (prompt.dataset.platforms || 'ios android').split(/\s+/);
        var key = prompt.dataset.dismissKey || 'default';
        if (!platforms.includes(platform) || isDismissed(key)) return;

        var link = prompt.querySelector('[data-app-prompt-link]');
        if (link) forwardLocation(link);

        var close = prompt.querySelector('[data-app-prompt-dismiss]');
        if (close) close.addEventListener('click', function () { dismiss(prompt); });
        prompt.hidden = false;
    });
})();
