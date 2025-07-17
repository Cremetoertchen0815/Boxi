// noinspection JSUnresolvedReference
const lightingModeOptions = [
    $('#overwrite-lighting-palette-row')[0],
    $('#overwrite-lighting-color-a-row')[0],
    $('#overwrite-lighting-color-b-row')[0],
    $('#overwrite-lighting-duration-row')[0],
    $('#overwrite-lighting-frequency-row')[0],
    $('#overwrite-lighting-brightness-row')[0],
    $('#overwrite-lighting-speed-row')[0],
    $('#overwrite-lighting-shift-row')[0]
];

const lightingModeSelector = $('#overwrite-lighting-mode')[0];

//On lighting override checkbox changed
$('#overwrite-lighting').on('change', async e => {
    const checked = e.target.checked;
    const contentPanel = $('#overwrite-lighting-content')[0];
    contentPanel.style.display = checked ? "block" : "none";

    if (checked) return;
    await sendLightingOverwrite(true, false);
});

//On animation override checkbox changed
$('#overwrite-animation').on('change', async e => {
    const checked = e.target.checked;
    const contentPanel = $('#overwrite-animation-content')[0];
    contentPanel.style.display = checked ? "block" : "none";

    if (checked) return;
    await sendAnimationOverwrite(true);
});

//On text override checkbox changed
$('#overwrite-text').on('change', async e => {
    const checked = e.target.checked;
    const contentPanel = $('#overwrite-text-content')[0];
    contentPanel.style.display = checked ? "block" : "none";

    if (checked) return;
    await sendTextsOverwrite(true);
});


//On lighting override mode changed
lightingModeSelector.onchange = _ => {
    //Select the options to be modified
    const index = parseInt(lightingModeSelector.selectedOptions[0].value);
    switch (index) {
        case 1: //Solid color
            setLightingModeOptions([1, 2]);
            break;
        case 2: //Fade to color
            setLightingModeOptions([1, 2, 3]);
            break;
        case 3: //Palette fade
            setLightingModeOptions([0, 3, 7]);
            break;
        case 4: //Palette switch
            setLightingModeOptions([0, 7]);
            break;
        case 5: //Brightness flash
            setLightingModeOptions([0, 5, 6, 7]);
            break;
        case 6: //Hue flash
            setLightingModeOptions([0, 6, 7]);
            break;
        case 7: //Strobe
            setLightingModeOptions([1, 4]);
            break;
        default: //Off
            setLightingModeOptions([]);
    }
};

function setLightingModeOptions(indices) {
    console.log(indices);
    for (const index of indices) {
        lightingModeOptions[index].style.display = "table-row";
    }

    for (let i = 0; i<lightingModeOptions.length; i++) {
        if (indices.some(v => v === i)) {
            continue;
        }

        lightingModeOptions[i].style.display = "none";
    }
}

$('#overwrite-lighting-apply-now').on('click', async _ => await sendLightingOverwrite(false, false));

$('#overwrite-lighting-apply-on-beat').on('click', async _ => await sendLightingOverwrite(false, true));

$('#overwrite-animation-apply').on('click', async _ => await sendAnimationOverwrite(false));

$('#overwrite-texts-apply').on('click', async _ => await sendTextsOverwrite(false));

async function sendLightingOverwrite(reset, onBeat) {
    const returnObj = {
        enable: !reset,
        onBeat: onBeat,
        mode: parseInt(lightingModeSelector.selectedOptions[0].value),
        colorA: getColorFromColorPicker($('#overwrite-lighting-color-a')[0]),
        colorB: getColorFromColorPicker($('#overwrite-lighting-color-b')[0]),
        paletteId: parseInt($('#overwrite-lighting-palette')[0].selectedOptions[0].value),
        duration: parseInt($('#overwrite-lighting-duration')[0].value),
        paletteShift: parseInt($('#overwrite-lighting-shift')[0].value),
        speed: parseInt($('#overwrite-lighting-speed')[0].value),
        targetBrightness: parseInt($('#overwrite-lighting-brightness')[0].value),
        frequency: parseInt($('#overwrite-lighting-frequency')[0].value)
    }

    await fetch(baseAddr + 'api/lighting/mode', {
        method: 'POST',
        body: JSON.stringify(returnObj),
    });
}

async function sendAnimationOverwrite(reset) {
    const animations = [];
    $('.overwrite-animation-id').each((_, animationSelector) =>
        animations.push({
            screen: parseInt(animationSelector.getAttribute('index')),
            animationId: parseInt(animationSelector.selectedOptions[0].value)
        }));

    const returnObj = {
        animations: animations,
        reset: reset,
        fadeoutSpeed: parseInt($('#overwrite-animation-fadeout')[0].value)
    }

    await fetch(baseAddr + 'api/screen/animation', {
        method: 'POST',
        body: JSON.stringify(returnObj),
    });
}

async function sendTextsOverwrite(reset) {
    const texts = [];
    $('.overwrite-text-id').each((_, animationSelector) =>
        texts.push({
            screen: parseInt(animationSelector.getAttribute('index')),
            text: reset ? " " : animationSelector.value
        }));

    const returnObj = { texts: texts }

    await fetch(baseAddr + 'api/screen/text', {
        method: 'POST',
        body: JSON.stringify(returnObj),
    });
}