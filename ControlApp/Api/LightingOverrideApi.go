package Api

import (
	"ControlApp/BoxiBus"
	"ControlApp/Lightshow"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	fadeDurationMsToCycles float64 = 500
)

type lightingInstructionSetColor struct {
	ColorDeviceA color
	ColorDeviceB color
}

type lightingInstructionFadeToColor struct {
	ColorDeviceA color
	ColorDeviceB color
	DurationMs   int
}

type lightingInstructionPaletteFade struct {
	PaletteId  uint32
	DurationMs int
}

func (fixture Fixture) HandleSetLightingOverrideAutoApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	fixture.Visuals.SetLightingOverwrite(nil)
}

func (fixture Fixture) HandleSetLightingOverrideOffApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	instruction := Lightshow.LightingInstruction{
		MessageBlock: BoxiBus.CreateLightingOff(false),
	}

	fixture.Visuals.SetLightingOverwrite(&instruction)
}

func (fixture Fixture) HandleSetLightingOverrideSetColorApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var data lightingInstructionSetColor

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid parameters. %s", err), http.StatusBadRequest)
		return
	}

	color1 := BoxiBus.Color{
		Red:         data.ColorDeviceA.R,
		Green:       data.ColorDeviceA.G,
		Blue:        data.ColorDeviceA.B,
		White:       data.ColorDeviceA.W,
		Amber:       data.ColorDeviceA.A,
		UltraViolet: data.ColorDeviceA.UV,
	}

	color2 := BoxiBus.Color{
		Red:         data.ColorDeviceB.R,
		Green:       data.ColorDeviceB.G,
		Blue:        data.ColorDeviceB.B,
		White:       data.ColorDeviceB.W,
		Amber:       data.ColorDeviceB.A,
		UltraViolet: data.ColorDeviceB.UV,
	}

	instruction := Lightshow.LightingInstruction{
		MessageBlock: BoxiBus.CreateLightingSetColor(color1, color2, false),
	}

	fixture.Visuals.SetLightingOverwrite(&instruction)
}

func (fixture Fixture) HandleSetLightingOverrideFadeToColorApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var data lightingInstructionFadeToColor

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid parameters. %s", err), http.StatusBadRequest)
		return
	}

	durationCycles := int(float64(data.DurationMs) * fadeDurationMsToCycles)
	if durationCycles <= 0 || durationCycles > 0xFFFF {
		http.Error(w, "Fade duration outside of range.", http.StatusBadRequest)
		return
	}

	color1 := BoxiBus.Color{
		Red:         data.ColorDeviceA.R,
		Green:       data.ColorDeviceA.G,
		Blue:        data.ColorDeviceA.B,
		White:       data.ColorDeviceA.W,
		Amber:       data.ColorDeviceA.A,
		UltraViolet: data.ColorDeviceA.UV,
	}

	color2 := BoxiBus.Color{
		Red:         data.ColorDeviceB.R,
		Green:       data.ColorDeviceB.G,
		Blue:        data.ColorDeviceB.B,
		White:       data.ColorDeviceB.W,
		Amber:       data.ColorDeviceB.A,
		UltraViolet: data.ColorDeviceB.UV,
	}

	instruction := Lightshow.LightingInstruction{
		MessageBlock: BoxiBus.CreateLightingFadeToColor(color1, color2, uint16(durationCycles), false),
	}

	fixture.Visuals.SetLightingOverwrite(&instruction)
}

func (fixture Fixture) HandleSetLightingOverridePaletteFadeApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var data lightingInstructionPaletteFade

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid parameters. %s", err), http.StatusBadRequest)
		return
	}

	fixture.Visuals.GetPalettes().GetById(data.PaletteId)

	durationCycles := int(float64(data.DurationMs) * fadeDurationMsToCycles)
	if durationCycles <= 0 || durationCycles > 0xFFFF {
		http.Error(w, "Fade duration outside of range.", http.StatusBadRequest)
		return
	}

	instruction := Lightshow.LightingInstruction{
		MessageBlock: BoxiBus.CreateLightingPaletteFade(color1, color2, uint16(durationCycles), false),
	}

	fixture.Visuals.SetLightingOverwrite(&instruction)
}
