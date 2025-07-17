// noinspection JSUnresolvedReference
const contentPanel = $('#palette-content')[0];
const countSelector = $('#palette-count')[0];
const nameInput = $('#palette-name')[0];
const itemSelection = $('#itemSelection')[0];
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
const moodCheckboxes = [
    $('#palette-mood-happy')[0],
    $('#palette-mood-moody')[0],
    $('#palette-mood-regular')[0],
    $('#palette-mood-party')[0]
];
let currentPalette = -1;

itemSelection.onchange = async () => {
    currentPalette = parseInt(itemSelection.selectedOptions[0].value);
    if (currentPalette < 0) {
        contentPanel.style.display = "none";
        return;
    }

    const paletteData = await getColors(currentPalette);
    countSelector.value = paletteData.colors.length;
    nameInput.value = paletteData.name;

    for (let i = 0; i < 4; i++) {
        moodCheckboxes[i].checked = false;
    }

    for (let i = 0; i < paletteData.moods.length; i++) {
        moodCheckboxes[paletteData.moods[i]].checked = true;
    }

    for (let i = 0; i < 8; i++) {
        if (i >= paletteData.colors.length) {
            colorPickerRows[i].style.display = "none";
            continue;
        }

        colorPickers[i].setAttribute('color', getColorString(paletteData.colors[i]));
        colorPickerRows[i].style.display = "block";
    }

    contentPanel.style.display = "block";
};

$('#palette-save').on('click', async _ => {
    await saveChange()

    if (currentPalette < 0) return;

    $('#itemSelection').children('option').each(function() {
        if (parseInt(this.value) === currentPalette) {
            this.text = nameInput.value;
        }
    })
});

countSelector.onchange = () => {
    const value = parseInt(countSelector.value);
    if (value < 1 || value > 8) return;

    for (let i = 0; i<8; i++) {
        const pickerRow = colorPickerRows[i];
        pickerRow.style.display = i < value ? "initial" : "none";
    }
};

$('#itemAddButton').on('click', async () => {
    const name = prompt("Please enter the name of the new palette:");
    if (name === null) return;

    const response = await createPalette(name);
    const id = response.id;

    const itemOption = document.createElement("option");
    itemOption.value = id;
    itemOption.innerText = name;
    itemSelection.appendChild(itemOption);
    itemSelection.value = id;

    await itemSelection.onchange();
});

$('#itemRemoveButton').on('click', async () => {
    const idToDelete = itemSelection.selectedOptions[0].value;
    const name = itemSelection.selectedOptions[0].innerText;
    if (idToDelete < 0) return;
    if (!confirm("Do you really want to delete palette '" + name + "'?")) return;

    await deletePalette(idToDelete);

    let childToDelete = null;
    const children = itemSelection.children;
    for (let i = 0; i < children.length; i++) {
        const c = children[i];
        if (c.value !== idToDelete) continue;
        childToDelete = c;
        break;
    }

    if (childToDelete != null) {
        childToDelete.remove();
    }

    itemSelection.value = -1;
    await itemSelection.onchange();
});

function getColorString(colors) {
    return [colors.R, colors.G, colors.B, colors.W, colors.A, colors.UV].join(',');
}

async function getColors(id) {
    const texts = [];
    $('.overwrite-text-id').each((_, animationSelector) =>
        texts.push({
            screen: parseInt(animationSelector.getAttribute('index')),
            text: reset ? " " : animationSelector.value
        }));

    const result = await fetch(baseAddr + 'api/palette?id=' + id, {
        method: 'GET'
    });

    return await result.json()
}

async function saveChange() {

    const moods = [];
    for (let i = 0; i<4; i++) {
        if (!moodCheckboxes[i].checked) continue;
        moods.push(i);
    }

    const colorCount = parseInt(countSelector.value);
    const color = [];
    for (let i = 0; i<colorCount; i++) {
        const picker = colorPickers[i];
        color.push(getColorFromColorPicker(picker));
    }

    const data = {
        id: currentPalette,
        name: nameInput.value,
        moods: moods,
        colors: color
    }

    await fetch(baseAddr + 'api/palette', {
        method: 'PUT',
        body: JSON.stringify(data)
    });
}

async function createPalette(name) {
    const data = {
        name: name,
    }

    const result = await fetch(baseAddr + 'api/palette', {
        method: 'POST',
        body: JSON.stringify(data)
    });

    return await result.json()
}

async function deletePalette(id) {
    await fetch(baseAddr + 'api/palette?id=' + id, {
        method: 'DELETE'
    });
}