$('#auto-config-save').on('click', async () => {
    const result = {
        strobeChance: parseInt($('#auto-config-strobe-chance')[0].value),
        hueShiftChance: parseInt($('#auto-config-hue-shift-chance')[0].value),
        fadeToColorDuration: parseInt($('#auto-config-duration-ftc')[0].value),
        paletteFadeDuration: parseInt($('#auto-config-duration-palette-fade')[0].value),
        brightnessFlashFadeSpeed: parseInt($('#auto-config-speed-brightness-flash')[0].value),
        hueFlashFadeSpeed: parseInt($('#auto-config-speed-hue-flash')[0].value),
        strobeFrequency: parseInt($('#auto-config-strobe-frequency')[0].value),
        brightnessFlashBrightness: parseInt($('#auto-config-target-brightness-flash')[0].value),
        hueFlashShift: parseInt($('#auto-config-color-span-hue-flash')[0].value),
        minTimeBetweenBeats: parseInt($('#auto-config-beat-cooldown-time')[0].value),
        timeBeforeLightingBoring: parseInt($('#auto-config-calm-lighting-boring')[0].value),
        timeBeforeAnimationBoring: parseInt($('#auto-config-calm-animation-boring')[0].value),
        timingRhythmicLighting: {
            minBeatsUntilSwitch: parseInt($('#auto-config-rhythmic-lighting-min-beats')[0].value),
            maxBeatsUntilSwitch: parseInt($('#auto-config-rhythmic-lighting-max-beats')[0].value),
            noBeatDeadTime: parseInt($('#auto-config-rhythmic-lighting-calmdown')[0].value),
        },
        timingRhythmicAnimations: {
            minBeatsUntilSwitch: parseInt($('#auto-config-rhythmic-animation-min-beats')[0].value),
            maxBeatsUntilSwitch: parseInt($('#auto-config-rhythmic-animation-max-beats')[0].value),
            noBeatDeadTime: parseInt($('#auto-config-rhythmic-animation-calmdown')[0].value),
        },
        timingFranticLighting: {
            minBeatsUntilSwitch: parseInt($('#auto-config-frantic-lighting-min-beats')[0].value),
            maxBeatsUntilSwitch: parseInt($('#auto-config-frantic-lighting-max-beats')[0].value),
            noBeatDeadTime: parseInt($('#auto-config-frantic-lighting-calmdown')[0].value),
        },
        timingFranticAnimations: {
            minBeatsUntilSwitch: parseInt($('#auto-config-frantic-animation-min-beats')[0].value),
            maxBeatsUntilSwitch: parseInt($('#auto-config-frantic-animation-max-beats')[0].value),
            noBeatDeadTime: parseInt($('#auto-config-frantic-animation-calmdown')[0].value),
        },
    };

    await fetch(baseAddr + 'api/config/advanced', {
        method: 'POST',
        body: JSON.stringify(result),
    });
});