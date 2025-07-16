
$('#animation-add-icon').on('click', e => {
   e.target.style.display = 'none';

   $('#animation-add-container')[0].classList.add('animation-add-container-active');
   $('#animation-add-dialog')[0].style.display = 'block';
});

$('#data-form').on('submit', async e => {
    e.preventDefault();

    const form = e.target;
    const url = form.action;
    const formData = new FormData(form);

    // UI Elements
    const status = $('#upload-status')[0];
    const errorBox = $('#upload-error')[0];

    // Show spinner and hide previous errors
    status.style.display = 'block';
    errorBox.style.display = 'none';
    errorBox.textContent = '';

    try {
        const response = await fetch(url, {
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            const errorText = await response.text();
            // noinspection ExceptionCaughtLocallyJS
            throw new Error(errorText || `Upload failed with status ${response.status}`);
        }

        // Success - reload the page
        location.reload();
    } catch (err) {
        // Show error message
        errorBox.textContent = err.message;
        errorBox.style.display = 'block';
    } finally {
        // Hide spinner whether success or failure
        status.style.display = 'none';
    }
});