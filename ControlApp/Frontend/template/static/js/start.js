$('#mood').change(moodChanged)

async function moodChanged(e) {
    const selected = e.target.selectedOptions[0];
    const value = selected.value;
    const response = await sendMoodChange(value);

    if (response != null) {
        alert("Error submitting mood change: " + response);
    }
}

async function moodChanged(e) {
    const selected = e.target.selectedOptions[0];
    const value = selected.value;
    const response = await sendMoodChange(value);

    if (response != null) {
        alert("Error submitting mood change: " + response);
    }
}

async function sendMoodChange(mood) {
    const response = await fetch('/api/config/mood?' + new URLSearchParams({
        value: mood
    }), {method: 'POST'});
    return response.status === 200 ? null : await response.text();
}