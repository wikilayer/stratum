// Copy-to-clipboard helpers. Two callsites today:
//
//   1. `.url-pill` rendered server-side with a `<button class="copy-btn">`
//      next to a `<code>` holding the value. The button stays visible
//      because the pill is itself the affordance ("here's a URL,
//      copy it"). data-label-copy / data-label-copied carry the
//      localised labels so this script never imports translation
//      strings.
//
//   2. Multi-line fenced code blocks (`<pre><code>…</code></pre>`).
//      The button is injected client-side and sits absolute top-right
//      over each block, fading in on hover. Labels come from
//      `<html data-label-copy=… data-label-copied=…>` so the server
//      controls the locale.
(function () {
    var docLabels = document.documentElement.dataset;
    var defaultCopy = docLabels.labelCopy || 'Copy';
    var defaultCopied = docLabels.labelCopied || 'Copied';

    function wire(btn, getValue, labelCopy, labelCopied) {
        btn.addEventListener('click', function () {
            var value = getValue();
            if (!navigator.clipboard || !value) return;
            navigator.clipboard.writeText(value).then(function () {
                btn.textContent = labelCopied;
                btn.classList.add('copied');
                setTimeout(function () {
                    btn.textContent = labelCopy;
                    btn.classList.remove('copied');
                }, 1500);
            });
        });
    }

    // (1) Existing url-pill buttons — rendered server-side.
    document.querySelectorAll('.url-pill .copy-btn').forEach(function (btn) {
        var pill = btn.closest('.url-pill');
        var code = pill && pill.querySelector('code');
        if (!code) return;
        var labelCopy = btn.dataset.labelCopy || btn.textContent;
        var labelCopied = btn.dataset.labelCopied || defaultCopied;
        wire(btn, function () { return code.textContent.trim(); }, labelCopy, labelCopied);
    });

    // (2) Fenced code blocks — inject a button per `<pre><code>` whose
    // body has at least one newline. Single-line fenced blocks read as
    // inline snippets in prose and don't need the affordance.
    document.querySelectorAll('pre > code').forEach(function (code) {
        if (!code.textContent.includes('\n')) return;
        var pre = code.parentElement;
        // Idempotent: don't double-inject if the script runs twice
        // (hot reload, multiple includes).
        if (pre.querySelector(':scope > .copy-btn')) return;
        var btn = document.createElement('button');
        btn.type = 'button';
        btn.className = 'copy-btn';
        btn.textContent = defaultCopy;
        pre.appendChild(btn);
        wire(btn, function () { return code.textContent.replace(/\n$/, ''); }, defaultCopy, defaultCopied);
    });
})();
