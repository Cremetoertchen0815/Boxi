package Api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type autoModeConfig struct {
	Mood                       int              `json:"mood"`
	AllowNsfw                  bool             `json:"nsfw"`
	StrobeChance               int              `json:"strobeChance"`
	HueShiftChance             int              `json:"hueShiftChance"`
	FadeToColorCycles          uint16           `json:"fadeToColorCycles"`
	PaletteFadeCycles          uint16           `json:"paletteFadeCycles"`
	FlashFadeoutSpeed          uint16           `json:"brightnessFlashFadeSpeed"`
	HueFlashFadeoutSpeed       uint16           `json:"hueFlashFadeSpeed"`
	StrobeFrequency            uint16           `json:"strobeFrequency"`
	FlashTargetBrightness      byte             `json:"brightnessFlashBrightness"`
	FlashHueShift              byte             `json:"hueFlashShift"`
	MinTimeBetweenBeatsSec     float64          `json:"minTimeBetweenBeats"`
	LightingCalmModeBoringSec  float64          `json:"timeBeforeLightingBoring"`  //How long it takes until calm lighting is boring
	AnimationCalmModeBoringSec float64          `json:"timeBeforeAnimationBoring"` //How long it takes until a calm animation is boring
	CalmLightingTiming         timingConstraint `json:"timingCalmLighting"`        //The timing constraints for calm lighting
	RhythmicLightingTiming     timingConstraint `json:"timingRhythmicLighting"`    //The timing constraints for rhythmic lighting
	FranticLightingTiming      timingConstraint `json:"timingFranticLighting"`     //The timing constraints for frantic lighting
	CalmAnimationsTiming       timingConstraint `json:"timingCalmAnimations"`      //The timing constraints for calm animations
	RhythmicAnimationsTiming   timingConstraint `json:"timingRhythmicAnimations"`  //The timing constraints for rhythmic animations
	FranticAnimationsTiming    timingConstraint `json:"timingFranticAnimations"`   //The timing constraints for calm animations
}

type timingConstraint struct {
	MinNumberOfBeats int           //The least number of beats before switching to the next mode.
	MaxNumberOfBeats int           //The most number of beats before switching to the next mode.
	NoBeatDeadTime   time.Duration //The duration since the last beat when forcibly switching to a calm mode.
}

func (fixture Fixture) HandleChangeAutoModeApi(w http.ResponseWriter, r *http.Request) {
	var data autoModeConfig

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid parameters. %s", err), http.StatusBadRequest)
		return
	}

	fixture.Visuals.MarkLightshowAsDirty()
}
