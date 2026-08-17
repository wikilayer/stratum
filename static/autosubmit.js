// A control that saves itself. Put data-autosubmit on a form and any
// change to a select, checkbox or radio inside it submits the form,
// so a one-field setting needs no Save button to press and no state
// to remember between choosing and saving.
//
// Text inputs are deliberately left out: every keystroke is a change,
// and a form that posts mid-word is worse than a button.

(function () {
    const forms = document.querySelectorAll('form[data-autosubmit]');
    if (forms.length === 0) return;

    forms.forEach(function (form) {
        form.addEventListener('change', function (e) {
            const el = e.target;
            if (!el.matches('select, input[type="checkbox"], input[type="radio"]')) return;
            // Disable the control while the page navigates, so a second
            // choice cannot race the first one's request.
            el.disabled = true;
            form.submit();
        });
    });
})();
