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
        let sent = false;
        form.addEventListener('change', function (e) {
            if (!e.target.matches('select, input[type="checkbox"], input[type="radio"]')) return;
            // A second change while the first request is in flight would
            // race it, and the loser decides what gets stored. The guard
            // is a flag rather than disabling the control: a disabled
            // field is not serialised, so disabling it here would post
            // the form without the very value that changed.
            if (sent) return;
            sent = true;
            form.submit();
        });
    });
})();
