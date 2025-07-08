
async function sendRegisterGame(name) {
    const response = await fetch('/api/quiz/round?' + new URLSearchParams({
        name: name
    }), {method: 'PUT'});
    return response.status === 200 ? await response.json() : null;
}