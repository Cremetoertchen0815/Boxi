package Lightshow

import (
	"encoding/json"
	"log"
	"os"
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

func loadConfiguration() AutoModeConfiguration {
	configFile, err := os.Open(autoModeConfigPath)

	var config AutoModeConfiguration
	if err != nil {
		log.Fatalf("Config file for auto mode could not be accessed! %s", err)
	}

	defer func(configFile *os.File) {
		_ = configFile.Close()
	}(configFile)

	jsonParser := json.NewDecoder(configFile)

	err = jsonParser.Decode(&config)
	if err != nil {
		log.Fatalf("Invalid JSON format of auto mode config file! %s", err)
	}

	return config
}

func storeConfiguration(config *AutoModeConfiguration) {
	configFile, err := os.OpenFile(autoModeConfigPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)

	if err != nil {
		log.Fatalf("Config file for auto mode could not be opened for writing! %s", err)
	}

	defer func(configFile *os.File) {
		_ = configFile.Close()
	}(configFile)

	jsonParser := json.NewEncoder(configFile)
	err = jsonParser.Encode(config)
	if err != nil {
		log.Fatalf("Configuration for auto mode could be JSON encoded! %s", err)
	}
}

// IsCalm returns whether the mood has exclusively calm character.
func (mood LightingMood) IsCalm() bool {
	return mood == Moody || mood == Happy
}
