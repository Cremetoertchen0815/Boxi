package Api

import (
	"ControlApp/Infrastructure"
	"ControlApp/Lightshow"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type AutoModeConfig struct {
	StrobeChance               int              `json:"strobeChance"`
	HueShiftChance             int              `json:"hueShiftChance"`
	FadeToColorMs              uint16           `json:"fadeToColorDuration"`
	PaletteFadeMs              uint16           `json:"paletteFadeDuration"`
	FlashFadeoutSpeed          uint16           `json:"brightnessFlashFadeSpeed"`
	HueFlashFadeoutSpeed       uint16           `json:"hueFlashFadeSpeed"`
	StrobeFrequency            uint16           `json:"strobeFrequency"`
	FlashTargetBrightness      byte             `json:"brightnessFlashBrightness"`
	FlashHueShift              byte             `json:"hueFlashShift"`
	MinTimeBetweenBeatsMs      uint16           `json:"minTimeBetweenBeats"`
	LightingCalmModeBoringSec  uint16           `json:"timeBeforeLightingBoring"`  //How long it takes until calm lighting is boring
	AnimationCalmModeBoringSec uint16           `json:"timeBeforeAnimationBoring"` //How long it takes until a calm animation is boring
	RhythmicLightingTiming     TimingConstraint `json:"timingRhythmicLighting"`    //The timing constraints for rhythmic lighting
	FranticLightingTiming      TimingConstraint `json:"timingFranticLighting"`     //The timing constraints for frantic lighting
	RhythmicAnimationsTiming   TimingConstraint `json:"timingRhythmicAnimations"`  //The timing constraints for rhythmic animations
	FranticAnimationsTiming    TimingConstraint `json:"timingFranticAnimations"`   //The timing constraints for calm animations
}

type TimingConstraint struct {
	MinNumberOfBeats  int     `json:"minBeatsUntilSwitch"` //The least number of beats before switching to the next mode.
	MaxNumberOfBeats  int     `json:"maxBeatsUntilSwitch"` //The most number of beats before switching to the next mode.
	NoBeatDeadTimeSec float64 `json:"noBeatDeadTime"`      //The duration since the last beat when forcibly switching to a calm mode.
}

func (fixture Fixture) HandleChangeAutoModeMoodApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	//Get animation ID
	var valueNr uint32
	valueNrStr := r.FormValue("value")
	if valueNrStr != "" {
		tempId, err := strconv.ParseInt(valueNrStr, 10, 8)
		if err != nil || tempId < 0 || tempId > 3 {
			http.Error(w, fmt.Sprintf("Error parsing animation ID. %s", err), http.StatusBadRequest)
			return
		}
		valueNr = uint32(tempId)
	}
	configuration := fixture.Data.Visuals.GetConfiguration()
	configuration.Mood = Lightshow.LightingMood(valueNr)
	fixture.Data.Visuals.StoreConfiguration(true)
}

func (fixture Fixture) HandleChangeAutoModeNsfwApi(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	//Get animation ID
	value := r.FormValue("value") == "true" || r.FormValue("value") == "1"
	configuration := fixture.Data.Visuals.GetConfiguration()
	configuration.AllowNsfw = value
	fixture.Data.Visuals.StoreConfiguration(true)
}

func (fixture Fixture) HandleChangeAutoModeConfigApi(w http.ResponseWriter, r *http.Request) {
	var data AutoModeConfig

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid parameters. %s", err), http.StatusBadRequest)
		return
	}

	configuration := fixture.Data.Visuals.GetConfiguration()
	configuration.StrobeChance = data.StrobeChance
	configuration.HueShiftChance = data.HueShiftChance
	configuration.FadeToColorCycles = uint16(float64(data.FadeToColorMs) * Infrastructure.FadeDurationMsToCycles)
	configuration.PaletteFadeCycles = uint16(float64(data.PaletteFadeMs) * Infrastructure.FadeDurationMsToCycles)
	configuration.FlashFadeoutSpeed = data.FlashFadeoutSpeed
	configuration.HueFlashFadeoutSpeed = data.HueFlashFadeoutSpeed
	configuration.StrobeFrequency = uint16(Infrastructure.StrobeFrequencyMultiplier / float64(data.StrobeFrequency))
	configuration.FlashTargetBrightness = byte(float64(data.FlashTargetBrightness) / 100 * 255)
	configuration.FlashHueShift = data.FlashHueShift
	configuration.MinTimeBetweenBeats = time.Duration(data.MinTimeBetweenBeatsMs) * time.Millisecond
	configuration.LightingCalmModeBoring = time.Duration(data.LightingCalmModeBoringSec) * time.Second
	configuration.AnimationCalmModeBoring = time.Duration(data.AnimationCalmModeBoringSec) * time.Second
	configuration.LightingModeTiming[Lightshow.Rhythmic] = getConstraint(data.RhythmicLightingTiming)
	configuration.LightingModeTiming[Lightshow.Frantic] = getConstraint(data.FranticLightingTiming)
	configuration.AnimationModeTiming[Lightshow.Rhythmic] = getConstraint(data.RhythmicAnimationsTiming)
	configuration.AnimationModeTiming[Lightshow.Frantic] = getConstraint(data.FranticAnimationsTiming)
	fixture.Data.Visuals.StoreConfiguration(false)
}

func getConstraint(constraint TimingConstraint) Lightshow.TimingConstraint {
	return Lightshow.TimingConstraint{
		MinNumberOfBeats: constraint.MinNumberOfBeats,
		MaxNumberOfBeats: constraint.MaxNumberOfBeats,
		NoBeatDeadTime:   time.Duration(float64(time.Second) * constraint.NoBeatDeadTimeSec),
	}
}
