package Lightshow

import (
	"log"
	"time"
)

type AutoModeConfiguration struct {
	Mood                    LightingMood
	AllowNsfw               bool
	StrobeChance            int
	HueShiftChance          int
	HueShiftMaxAmount       int
	FadeToColorCycles       uint16 //How slow is the “FadeToColor” mode operating at
	PaletteFadeCycles       uint16 //How slow is the “FadeToColor” mode operating at
	FlashFadeoutSpeed       uint16
	HueFlashFadeoutSpeed    uint16
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

type LightingMood uint16

const (
	Happy LightingMood = iota
	Moody
	Regular
	Party
)

const autoModeConfigPath = "Configuration/auto_mode.json"

const autoModeConfigBackupPath = "Configuration/auto_mode_backup.json"

func load() AutoModeConfiguration {
	config, err := loadConfiguration[AutoModeConfiguration](autoModeConfigPath)

	if err != nil {
		config, err = loadConfiguration[AutoModeConfiguration](autoModeConfigBackupPath)
	}

	if err != nil {
		log.Fatalf("Config file for auto mode could not be accessed! %s", err)
	}

	return config
}

func (config *AutoModeConfiguration) Store() {
	storeConfiguration(config, autoModeConfigPath, autoModeConfigBackupPath)
}

// IsCalm returns whether the mood has exclusively calm character.
func (mood LightingMood) IsCalm() bool {
	return mood == Moody || mood == Happy
}
