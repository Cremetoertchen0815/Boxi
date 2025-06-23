package Lightshow

import (
	"time"
)

type AutoModeConfiguration struct {
	Mood                    LightingMood
	StrobeChance            int
	HueShiftChance          int
	HueShiftMaxAmount       int
	FadeToColorCycles       uint16 //How slow is the “FadeToColor” mode operating at
	PaletteFadeCycles       uint16 //How slow is the “FadeToColor” mode operating at
	FlashFadeoutSpeed       uint16
	StrobeFrequency         uint16
	StrobeRolloff           byte
	FlashTargetBrightness   byte
	FlashHueShift           byte
	MinTimeBetweenBeats     time.Duration
	LightingCalmModeBoring  time.Duration                      //How long it takes until a calm animation is boring
	AnimationCalmModeBoring time.Duration                      //How long it takes until a calm animation is boring
	LightingModeTiming      map[ModeCharacter]TimingConstraint //The timing constraints for lighting of any character
	AnimationModeTiming     map[ModeCharacter]TimingConstraint //The timing constraints for animations of any character
}

type TimingConstraint struct {
	MinNumberOfBeats int           //The least number of beats before switching to the next mode.
	MaxNumberOfBeats int           //The most number of beats before switching to the next mode.
	NoBeatDeadTime   time.Duration //The duration since the last beat when forcibly switching to a calm mode.
}

type LightingMood uint8

const (
	Happy LightingMood = iota
	Moody
	Regular
	Party
)

func loadConfiguration() AutoModeConfiguration {
	return AutoModeConfiguration{
		Mood:                    Party,
		StrobeChance:            4,
		HueShiftChance:          3,
		HueShiftMaxAmount:       3,
		FadeToColorCycles:       700,
		PaletteFadeCycles:       500,
		StrobeFrequency:         2,
		FlashFadeoutSpeed:       28,
		FlashTargetBrightness:   20,
		FlashHueShift:           1,
		MinTimeBetweenBeats:     300 * time.Millisecond,
		LightingCalmModeBoring:  20 * time.Second,
		AnimationCalmModeBoring: 30 * time.Second,
		LightingModeTiming: map[ModeCharacter]TimingConstraint{
			Calm:     {MinNumberOfBeats: 32, MaxNumberOfBeats: 128, NoBeatDeadTime: 4 * time.Second},
			Rhythmic: {MinNumberOfBeats: 16, MaxNumberOfBeats: 64, NoBeatDeadTime: 3 * time.Second},
			Frantic:  {MinNumberOfBeats: 1, MaxNumberOfBeats: 4, NoBeatDeadTime: 1 * time.Second},
		},
		AnimationModeTiming: map[ModeCharacter]TimingConstraint{
			Calm:     {MinNumberOfBeats: 32, MaxNumberOfBeats: 64, NoBeatDeadTime: 6 * time.Second},
			Rhythmic: {MinNumberOfBeats: 16, MaxNumberOfBeats: 48, NoBeatDeadTime: 2 * time.Second},
			Frantic:  {MinNumberOfBeats: 8, MaxNumberOfBeats: 16, NoBeatDeadTime: 2 * time.Second},
		},
	}
}

// IsCalm returns whether the mood has exclusively calm character.
func (mood LightingMood) IsCalm() bool {
	return mood == Moody || mood == Happy
}
