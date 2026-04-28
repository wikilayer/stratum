// Wires every "Copy" button on the page. Each button is expected to
// live inside a `.url-pill` next to a `<code>` element holding the
// value to copy. data-label-copy / data-label-copied carry the
// localised labels so this script never imports translation strings.
(function () {
    var buttons = document.querySelectorAll('.copy-btn');
    buttons.forEach(function (btn) {
        var pill = btn.closest('.url-pill');
        if (!pill) return;
        var code = pill.querySelector('code');
        if (!code) return;

        var labelCopy = btn.dataset.labelCopy || btn.textContent;
        var labelCopied = btn.dataset.labelCopied || 'Copied';

        btn.addEventListener('click', function () {
            var value = code.textContent.trim();
            if (!navigator.clipboard) return;
            navigator.clipboard.writeText(value).then(function () {
                btn.textContent = labelCopied;
                btn.classList.add('copied');
                setTimeout(function () {
                    btn.textContent = labelCopy;
                    btn.classList.remove('copied');
                }, 1500);
            });
        });
    });
})();
