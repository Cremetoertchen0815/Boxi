package Api

import (
	"ControlApp/Display"
	"ControlApp/Lightshow"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type ScreenOverrideAnimationProperties struct {
	Animations   []ScreenOverrideAnimationInstance `json:"animations"`
	FadeoutSpeed int                               `json:"fadeoutSpeed"`
	ResetScreens bool                              `json:"reset"`
}

type ScreenOverrideAnimationInstance struct {
	ScreenIndex Display.ServerDisplay `json:"screen"`
	AnimationId Display.AnimationId   `json:"animationId"`
}

type ScreenOverrideTextProperties struct {
	Texts []ScreenTextInstance `json:"texts"`
}

type ScreenTextInstance struct {
	ScreenIndex Display.ServerDisplay `json:"screen"`
	Text        string                `json:"text"`
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
	for _, index := range fixture.Data.Hardware.GetConnectedDisplays() {
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

	var data ScreenOverrideAnimationProperties

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid parameters. %s", err), http.StatusBadRequest)
		return
	}

	if data.ResetScreens {
		fixture.Data.Visuals.SetAnimationsOverwrite(nil)
		fixture.Data.Visuals.MarkLightshowAsDirty()
		return
	}

	var aniInstr []Lightshow.AnimationInstruction

	for _, animation := range data.Animations {
		if animation.ScreenIndex < Display.Boxi1D1 || animation.ScreenIndex > Display.Boxi2D2 {
			http.Error(w, fmt.Sprintf("Screen ID is out of bound. %s", err), http.StatusBadRequest)
			return
		}

		if animation.AnimationId == Display.None {
			continue
		}

		exists, animationObj := fixture.Data.Visuals.GetAnimations().GetById(animation.AnimationId)
		if !exists {
			http.Error(w, fmt.Sprintf("Animation can't be found. %s", err), http.StatusBadRequest)
			return
		}

		aniInstr = append(aniInstr, Lightshow.AnimationInstruction{Animation: animationObj.Id, Displays: []Display.ServerDisplay{animation.ScreenIndex}})
	}

	instr := Lightshow.AnimationsInstruction{Animations: aniInstr, Character: Lightshow.Unknown, BlinkSpeed: uint16(data.FadeoutSpeed)}
	fixture.Data.Visuals.SetAnimationsOverwrite(&instr)
	fixture.Data.Visuals.MarkLightshowAsDirty()
}

func (fixture Fixture) HandleSetScreenOverrideTextSetApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var data ScreenOverrideTextProperties

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid parameters. %s", err), http.StatusBadRequest)
		return
	}

	var textInstr []Lightshow.TextInstruction

	for _, text := range data.Texts {
		if text.ScreenIndex < Display.Boxi1D1 || text.ScreenIndex > Display.Boxi2D2 {
			http.Error(w, fmt.Sprintf("Screen ID is out of bound. %s", err), http.StatusBadRequest)
			return
		}
		textContent := text.Text
		if strings.TrimSpace(textContent) == "" {
			textContent = " "
		}

		textInstr = append(textInstr, Lightshow.TextInstruction{Text: text.Text, Displays: []Display.ServerDisplay{text.ScreenIndex}})
	}

	fixture.Data.Visuals.SetTexts(textInstr)
}

func (fixture Fixture) HandleSetScreenOverrideBrightnessLevelApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var tempId float64
	moodNrStr := r.FormValue("value")
	if moodNrStr != "" {
		tmp, err := strconv.ParseInt(moodNrStr, 10, 32)
		if err != nil || tempId < 0 {
			http.Error(w, "Error parsing mood.", http.StatusBadRequest)
			return
		}
		tempId = float64(tmp) / 100.0
	}

	fixture.Data.Visuals.SetBrightness(tempId)
}
