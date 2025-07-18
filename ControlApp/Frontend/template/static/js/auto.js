$('#auto-config-save').on('click', async () => {
    const result = {
        strobeChance: parseInt($('#auto-config-strobe-chance').value),
        hueShiftChance: parseInt($('#auto-config-hue-shift-chance').value),
        fadeToColorDuration: parseInt($('#auto-config-duration-ftc').value),
        paletteFadeDuration: parseInt($('#auto-config-duration-palette-fade').value),
        brightnessFlashFadeSpeed: parseInt($('#auto-config-speed-brightness-flash').value),
        hueFlashFadeSpeed: parseInt($('#auto-config-speed-hue-flash').value),
        strobeFrequency: parseInt($('#auto-config-strobe-frequency').value),
        brightnessFlashBrightness: parseInt($('#auto-config-target-brightness-flash').value),
        hueFlashShift: parseInt($('#auto-config-color-span-hue-flash').value),
        minTimeBetweenBeats: parseInt($('#auto-config-beat-cooldown-time').value),
        timeBeforeLightingBoring: parseInt($('#auto-config-calm-lighting-boring').value),
        timeBeforeAnimationBoring: parseInt($('#auto-config-calm-animation-boring').value),
        timingRhythmicLighting: {
            MinNumberOfBeats: parseInt($('#auto-config-rhythmic-lighting-min-beats').value),
            MaxNumberOfBeats: parseInt($('#auto-config-rhythmic-lighting-max-beats').value),
            NoBeatDeadTimeSec: parseInt($('#auto-config-rhythmic-lighting-calmdown').value),
        },
        timingFranticLighting: {
            MinNumberOfBeats: parseInt($('#auto-config-rhythmic-animation-min-beats').value),
            MaxNumberOfBeats: parseInt($('#auto-config-rhythmic-animation-max-beats').value),
            NoBeatDeadTimeSec: parseInt($('#auto-config-rhythmic-animation-calmdown').value),
        },
        timingRhythmicAnimations: {
            MinNumberOfBeats: parseInt($('#auto-config-frantic-lighting-min-beats').value),
            MaxNumberOfBeats: parseInt($('#auto-config-frantic-lighting-max-beats').value),
            NoBeatDeadTimeSec: parseInt($('#auto-config-frantic-lighting-calmdown').value),
        },
        timingFranticAnimations: {
            MinNumberOfBeats: parseInt($('#auto-config-frantic-animation-min-beats').value),
            MaxNumberOfBeats: parseInt($('#auto-config-frantic-animation-max-beats').value),
            NoBeatDeadTimeSec: parseInt($('#auto-config-frantic-animation-calmdown').value),
        },
    };

    await fetch(baseAddr + 'api/config/advanced', {
        method: 'POST',
        body: JSON.stringify(result),
    });
});