// noinspection JSUnresolvedReference
const baseAddr = "http://localhost:8080/"

const contentPanel = $('#palette-content')[0];
const countSelector = $('#palette-count')[0];
const colorPickers = [
    $('#palette-color-a')[0],
    $('#palette-color-b')[0],
    $('#palette-color-c')[0],
    $('#palette-color-d')[0],
    $('#palette-color-e')[0],
    $('#palette-color-f')[0],
    $('#palette-color-g')[0],
    $('#palette-color-h')[0],
];
const colorPickerRows = [
    $('#palette-color-a-row')[0],
    $('#palette-color-b-row')[0],
    $('#palette-color-c-row')[0],
    $('#palette-color-d-row')[0],
    $('#palette-color-e-row')[0],
    $('#palette-color-f-row')[0],
    $('#palette-color-g-row')[0],
    $('#palette-color-h-row')[0],
];
let currentPalette = -1;

$('#itemSelection').on('change', async e => {
    currentPalette = parseInt(e.target.selectedOptions[0].value);
    if (currentPalette < 0) {
        contentPanel.style.display = "none";
        return;
    }

    const paletteData = await getColors(currentPalette);
    console.log(paletteData);
    countSelector.value = paletteData.colors.length;

    for (let i = 0; i < 8; i++) {
        if (i >= paletteData.colors.length) {
            colorPickerRows[i].style.display = "none";
            continue;
        }

        colorPickers[i].setAttribute('color', getColorString(paletteData.colors[i]));
        colorPickerRows[i].style.display = "block";
    }

    contentPanel.style.display = "block";
});

function getColorString(colors) {
    console.log(colors);
    return [colors.R, colors.G, colors.B, colors.W, colors.A, colors.UV].join(',');
}

async function getColors(id) {
    const texts = [];
    $('.overwrite-text-id').each((_, animationSelector) =>
        texts.push({
            screen: parseInt(animationSelector.getAttribute('index')),
            text: reset ? " " : animationSelector.value
        }));

    const returnObj = { texts: texts }

    const result = await fetch(baseAddr + 'api/palette?id=' + id, {
        method: 'GET'
    });

    return await result.json()
}