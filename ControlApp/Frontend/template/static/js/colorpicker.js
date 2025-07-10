const channels = ['R', 'G', 'B', 'W', 'A', 'U'];
const values = [{ R: 0, G: 0, B: 0, W: 0, A: 0, U: 0 }, { R: 0, G: 0, B: 0, W: 0, A: 0, U: 0 }];

const targets = {
    W: [255, 255, 255],
    A: [255, 191, 0],
    U: [148, 0, 211]
};
$(".colorpicker").each(function(index, sliderContainer) {
    const colorPreview = document.createElement('div');
    colorPreview.className = "preview";
    sliderContainer.appendChild(colorPreview);

    // Parse the “color” attribute
    let initialValues = getColorArrayFromPicker(sliderContainer);

    // Initialize values
    channels.forEach((channel, i) => {
        values[index][channel] = initialValues[i];
    });

    // Build sliders
    channels.forEach(channel => {
        const column = document.createElement('div');
        column.className = 'slider-column';

        const slider = document.createElement('input');
        slider.type = 'range';
        slider.min = 0;
        slider.max = 255;
        slider.value = values[index][channel];
        slider.id = 'slider_' + channel;

        const label = document.createElement('div');
        label.className = 'channel-label';
        label.textContent = channel;

        const value = document.createElement('div');
        value.className = 'value-label';
        value.textContent = slider.value;
        value.id = 'value_' + channel;

        // Update logic
        slider.addEventListener('input', () => {
            values[index][channel] = parseInt(slider.value);
            value.textContent = slider.value;
            updatePreview(index, colorPreview);

            // Update the "color" attribute on this container
            const updatedColor = channels.map(c => values[index][c]).join(',');
            sliderContainer.setAttribute("color", updatedColor);
        });

        column.appendChild(label);
        column.appendChild(slider);
        column.appendChild(value);
        sliderContainer.appendChild(column);
    });

    // Initial preview render
    updatePreview(index, colorPreview);
});

function lerp(a, b, t) {
    return a + (b - a) * t;
}

function updatePreview(index, obj) {
    // Start with raw RGB base
    let r = values[index].R;
    let g = values[index].G;
    let b = values[index].B;

    // LERP each target into RGB color (W, A, U)
    const blend = { r, g, b };

    // Each max value (255) equals 50% blend
    const lerpFactor = (val) => val / 255 * 0.5;

    ['W', 'A', 'U'].forEach(channel => {
        const t = lerpFactor(values[index][channel]);
        const [tr, tg, tb] = targets[channel];

        blend.r = lerp(blend.r, tr, t);
        blend.g = lerp(blend.g, tg, t);
        blend.b = lerp(blend.b, tb, t);
    });

    // Clamp & round
    const rFinal = Math.round(Math.min(255, blend.r));
    const gFinal = Math.round(Math.min(255, blend.g));
    const bFinal = Math.round(Math.min(255, blend.b));

    obj.style.backgroundColor = `rgb(${rFinal}, ${gFinal}, ${bFinal})`;
}

function getColorArrayFromPicker(picker) {
    let initialValues = Array(6).fill(0);
    const attr = picker.getAttribute("color");
    if (attr) {
        const parts = attr.split(",").map(v => parseInt(v.trim()));
        if (parts.length === 6 && parts.every(n => !isNaN(n))) {
            initialValues = parts;
        }
    }

    return initialValues;
}

function getColorFromColorPicker(picker) {
    const initialValues = getColorArrayFromPicker(picker);
    return {
        R: initialValues[0],
        G: initialValues[1],
        B: initialValues[2],
        W: initialValues[3],
        A: initialValues[4],
        UV: initialValues[5],
    }
}