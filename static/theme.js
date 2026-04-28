// Theme toggle: writes a cookie and applies the change immediately
// without a page reload. Server-side reading of the cookie is what
// prevents flash-of-wrong-theme on subsequent navigations.

(function () {
    const buttons = document.querySelectorAll('[data-theme-set]');
    if (buttons.length === 0) return;

    function setTheme(value) {
        if (value === 'auto') {
            document.cookie = 'theme=; max-age=0; path=/; samesite=lax';
            document.documentElement.removeAttribute('data-theme');
        } else {
            document.cookie = 'theme=' + value + '; max-age=31536000; path=/; samesite=lax';
            document.documentElement.setAttribute('data-theme', value);
        }
        buttons.forEach(function (b) {
            b.setAttribute('aria-pressed', b.dataset.themeSet === value ? 'true' : 'false');
        });
    }

    buttons.forEach(function (b) {
        b.addEventListener('click', function () {
            setTheme(b.dataset.themeSet);
        });
    });
})();
