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
$('#overwrite-lighting').on('change', e => {
    const checked = e.target.checked;
    const contentPanel = $('#overwrite-lighting-content')[0];
    contentPanel.style.display = checked ? "block" : "none";

    //TODO: Send overwrite reset when unchecked
});

//On animation override checkbox changed
$('#overwrite-animation').on('change', e => {
    const checked = e.target.checked;
    const contentPanel = $('#overwrite-animation-content')[0];
    contentPanel.style.display = checked ? "block" : "none";

    //TODO: Send overwrite reset when unchecked
});

//On text override checkbox changed
$('#overwrite-text').on('change', e => {
    const checked = e.target.checked;
    const contentPanel = $('#overwrite-text-content')[0];
    contentPanel.style.display = checked ? "block" : "none";

    //TODO: Send overwrite reset when unchecked
});


//On lighting override mode changed
lightingModeSelector.on('change', e => {
    //Select the options to be modified
    const index = parseInt(e.target.selectedOptions[0].value);
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
});

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
