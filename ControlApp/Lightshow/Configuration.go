package Lightshow

import (
	"time"
)

type AutoModeConfiguration struct {
	Mood                    LightingMood
	StrobeChance            int
	HueShiftChance          int
	HueShiftMaxAmount       int
	FadeToColorCycles       uint16 //How slow is the "FadeToColor" mode operating at
	PaletteFadeCycles       uint16 //How slow is the "FadeToColor" mode operating at
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
	return AutoModeConfiguration{}
}

// IsCalm returns whether the mood has exclusively calm character.
func (mood LightingMood) IsCalm() bool {
	return mood == Moody || mood == Happy
}
