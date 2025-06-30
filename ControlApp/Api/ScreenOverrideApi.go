package Api

import (
	"ControlApp/Display"
	"ControlApp/Lightshow"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type screenOverrideAnimationProperties struct {
	Animations   []screenOverrideAnimationInstance `json:"animation"`
	FadeoutSpeed int                               `json:"fadeoutSpeed"`
	ResetScreens bool                              `json:"reset"`
}

type screenOverrideAnimationInstance struct {
	ScreenIndices []int  `json:"screen"`
	AnimationId   uint32 `json:"animationId"`
}

type screenOverrideTextProperties struct {
	Texts []screenTextInstance `json:"texts"`
}

type screenTextInstance struct {
	ScreenIndices []int  `json:"screen"`
	Text          string `json:"text"`
}

type fetchResult struct {
	ConnectedDisplays []int
}

func (fixture Fixture) HandleScreensConnectedApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	indices := make([]int, 0)
	for _, index := range fixture.Hardware.GetConnectedDisplays() {
		indices = append(indices, int(index))
	}

	//Encode data
	if err := json.NewEncoder(w).Encode(fetchResult{indices}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (fixture Fixture) HandleSetScreenOverrideAnimationSetApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var data screenOverrideAnimationProperties

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid parameters. %s", err), http.StatusBadRequest)
		return
	}

	if data.ResetScreens {
		fixture.Visuals.SetAnimationsOverwrite(nil)
		return
	}

	var aniInstr []Lightshow.AnimationInstruction

	for _, animation := range data.Animations {
		var indices []Display.ServerDisplay
		for _, index := range animation.ScreenIndices {
			if index < 0 || index > 3 {
				http.Error(w, fmt.Sprintf("Screen ID is out of bound. %s", err), http.StatusBadRequest)
				return
			}

			indices = append(indices, Display.ServerDisplay(index))
		}

		exists, animationObj := fixture.Visuals.GetAnimations().GetById(Display.AnimationId(animation.AnimationId))
		if !exists {
			http.Error(w, fmt.Sprintf("Animation can't be found. %s", err), http.StatusBadRequest)
			return
		}

		aniInstr = append(aniInstr, Lightshow.AnimationInstruction{Animation: animationObj.Id, Displays: indices})
	}

	instr := Lightshow.AnimationsInstruction{Animations: aniInstr, Character: Lightshow.Unknown, BlinkSpeed: uint16(data.FadeoutSpeed)}
	fixture.Visuals.SetAnimationsOverwrite(&instr)
}

func (fixture Fixture) HandleSetScreenOverrideTextSetApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var data screenOverrideTextProperties

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid parameters. %s", err), http.StatusBadRequest)
		return
	}

	var textInstr []Lightshow.TextInstruction

	for _, text := range data.Texts {
		var indices []Display.ServerDisplay
		for _, index := range text.ScreenIndices {
			if index < 0 || index > 3 {
				http.Error(w, fmt.Sprintf("Screen ID is out of bound. %s", err), http.StatusBadRequest)
				return
			}

			indices = append(indices, Display.ServerDisplay(index))
		}

		textInstr = append(textInstr, Lightshow.TextInstruction{Text: text.Text, Displays: indices})
	}

	fixture.Visuals.SetTexts(textInstr)
}

func (fixture Fixture) HandleSetScreenOverrideBrightnessLevelApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var tempId float64
	moodNrStr := r.FormValue("value")
	if moodNrStr != "" {
		tmp, err := strconv.ParseFloat(moodNrStr, 64)
		if err != nil || tempId < 0 {
			http.Error(w, "Error parsing mood.", http.StatusBadRequest)
			return
		}
		tempId = tmp
	}

	fixture.Visuals.SetBrightness(tempId)
}
