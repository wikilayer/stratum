// Modal helper. Wires:
//   <button data-modal-open="ID">    →  document.getElementById(ID).showModal()
//   <button data-modal-close>        →  closes the nearest enclosing <dialog>
//
// ESC and the close-button on the dialog header are native; this just
// adds open-by-id and click-on-backdrop-to-close. Markup convention:
//
//   <dialog id="x" class="modal">
//     <header class="modal-header">
//       <p class="modal-title">Title</p>
//       <button class="modal-close" data-modal-close aria-label="Close">×</button>
//     </header>
//     <div class="modal-body">…</div>
//     <footer class="modal-footer">
//       <button class="button" data-modal-close>Cancel</button>
//       <button class="button button-primary">OK</button>
//     </footer>
//   </dialog>

(function () {
    document.addEventListener('click', function (ev) {
        var openTrigger = ev.target.closest('[data-modal-open]');
        if (openTrigger) {
            ev.preventDefault();
            var dlg = document.getElementById(openTrigger.dataset.modalOpen);
            if (dlg && typeof dlg.showModal === 'function') dlg.showModal();
            return;
        }
        var closeTrigger = ev.target.closest('[data-modal-close]');
        if (closeTrigger) {
            ev.preventDefault();
            var dlg = closeTrigger.closest('dialog');
            if (dlg) dlg.close();
            return;
        }
    });

    // Click on the backdrop closes the dialog. The dialog element
    // catches the click and target === the dialog itself when the
    // user clicks outside the visible content rect.
    document.addEventListener('click', function (ev) {
        var dlg = ev.target;
        if (dlg.tagName !== 'DIALOG') return;
        var r = dlg.getBoundingClientRect();
        var inside =
            r.top  <= ev.clientY && ev.clientY <= r.top  + r.height &&
            r.left <= ev.clientX && ev.clientX <= r.left + r.width;
        if (!inside) dlg.close();
    });
})();
