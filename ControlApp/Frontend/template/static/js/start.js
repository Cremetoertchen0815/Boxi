$('#mood').change(moodChanged)
$('#nsfw').change(nsfwChanged)
$('#brightness').change(brightnessChanged)

async function moodChanged(e) {
    const selected = e.target.selectedOptions[0];
    const value = selected.value;
    const response = await sendMoodChange(value);

    if (response != null) {
        alert("Error submitting mood change: " + response);
    }
}

async function nsfwChanged(e) {
    const value = e.target.checked;
    const response = await sendNsfwChange(value);

    if (response != null) {
        alert("Error submitting NSFW change: " + response);
    }
}

async function brightnessChanged(e) {
    const value = Math.min(Math.max(parseInt(e.target.value), 0), 100);
    const response = await sendBrightnessChange(value);

    if (response != null) {
        alert("Error submitting brightness change: " + response);
    }
}

async function sendMoodChange(mood) {
    const response = await fetch('/api/config/mood?' + new URLSearchParams({
        value: mood
    }), {method: 'POST'});
    return response.status === 200 ? null : await response.text();
}

async function sendNsfwChange(allowNsfw) {
    const response = await fetch('/api/config/nsfw?' + new URLSearchParams({
        value: allowNsfw ? "true" : "false"
    }), {method: 'POST'});
    return response.status === 200 ? null : await response.text();
}

async function sendBrightnessChange(value) {
    const response = await fetch('/api/screen/brightness?' + new URLSearchParams({
        value: value
    }), {method: 'POST'});
    return response.status === 200 ? null : await response.text();
}