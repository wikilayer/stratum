(function () {
    var animations = new WeakMap();

    function milliseconds(value) {
        var number = parseFloat(value);
        if (!Number.isFinite(number)) return 120;
        return value.trim().endsWith('ms') ? number : number * 1000;
    }

    function setHidden(branch, hidden) {
        var running = animations.get(branch);
        if (running) running.cancel();

        if (matchMedia('(prefers-reduced-motion: reduce)').matches || !branch.animate) {
            branch.hidden = hidden;
            return;
        }

        if (!hidden) branch.hidden = false;
        branch.dataset.treeAnimating = '';

        var style = getComputedStyle(branch);
        var height = branch.scrollHeight + 'px';
        var frames = hidden
            ? [{ height: height, opacity: 1 }, { height: '0', opacity: 0 }]
            : [{ height: '0', opacity: 0 }, { height: height, opacity: 1 }];
        var animation = branch.animate(frames, {
            duration: milliseconds(style.getPropertyValue('--duration-fast')),
            easing: style.getPropertyValue('--easing-out').trim(),
        });
        animations.set(branch, animation);

        animation.finished.then(function () {
            if (animations.get(branch) !== animation) return;
            animations.delete(branch);
            delete branch.dataset.treeAnimating;
            if (hidden) branch.hidden = true;
        }).catch(function () {});
    }

    document.addEventListener('click', function (event) {
        var toggle = event.target.closest('.tree-nav-toggle');
        if (!toggle) return;

        var branch = document.getElementById(toggle.getAttribute('aria-controls'));
        if (!branch) return;

        var expanded = toggle.getAttribute('aria-expanded') === 'true';
        toggle.setAttribute('aria-expanded', String(!expanded));
        var label = toggle.getAttribute(expanded ? 'data-label-expand' : 'data-label-collapse');
        if (label) toggle.setAttribute('aria-label', label);
        setHidden(branch, expanded);
    });
})();
