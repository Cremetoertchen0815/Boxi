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
	ColorDeviceA color `json:"colorA"`
	ColorDeviceB color `json:"colorB"`
}

type lightingInstructionFadeToColor struct {
	ColorDeviceA color `json:"colorA"`
	ColorDeviceB color `json:"colorB"`
	DurationMs   int   `json:"duration"`
}

type lightingInstructionPaletteFade struct {
	PaletteId    uint32 `json:"paletteId"`
	DurationMs   int    `json:"duration"`
	PaletteShift int    `json:"paletteShift"`
}

type lightingInstructionPaletteSwitch struct {
	PaletteId    uint32 `json:"paletteId"`
	PaletteShift int    `json:"paletteShift"`
}

type lightingInstructionPaletteBrightnessFlash struct {
	PaletteId        uint32 `json:"paletteId"`
	TargetBrightness byte   `json:"targetBrightness"`
	DurationMs       int    `json:"duration"`
	PaletteShift     int    `json:"paletteShift"`
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

	if !isColorValid(data.ColorDeviceA) {
		http.Error(w, "Color A is invalid.", http.StatusBadRequest)
		return
	}

	if !isColorValid(data.ColorDeviceB) {
		http.Error(w, "Color B is invalid.", http.StatusBadRequest)
		return
	}

	color1 := BoxiBus.Color{
		Red:         byte(data.ColorDeviceA.R),
		Green:       byte(data.ColorDeviceA.G),
		Blue:        byte(data.ColorDeviceA.B),
		White:       byte(data.ColorDeviceA.W),
		Amber:       byte(data.ColorDeviceA.A),
		UltraViolet: byte(data.ColorDeviceA.UV),
	}

	color2 := BoxiBus.Color{
		Red:         byte(data.ColorDeviceB.R),
		Green:       byte(data.ColorDeviceB.G),
		Blue:        byte(data.ColorDeviceB.B),
		White:       byte(data.ColorDeviceB.W),
		Amber:       byte(data.ColorDeviceB.A),
		UltraViolet: byte(data.ColorDeviceB.UV),
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

	if !isColorValid(data.ColorDeviceA) {
		http.Error(w, "Color A is invalid.", http.StatusBadRequest)
		return
	}

	if !isColorValid(data.ColorDeviceB) {
		http.Error(w, "Color B is invalid.", http.StatusBadRequest)
		return
	}

	color1 := BoxiBus.Color{
		Red:         byte(data.ColorDeviceA.R),
		Green:       byte(data.ColorDeviceA.G),
		Blue:        byte(data.ColorDeviceA.B),
		White:       byte(data.ColorDeviceA.W),
		Amber:       byte(data.ColorDeviceA.A),
		UltraViolet: byte(data.ColorDeviceA.UV),
	}

	color2 := BoxiBus.Color{
		Red:         byte(data.ColorDeviceB.R),
		Green:       byte(data.ColorDeviceB.G),
		Blue:        byte(data.ColorDeviceB.B),
		White:       byte(data.ColorDeviceB.W),
		Amber:       byte(data.ColorDeviceB.A),
		UltraViolet: byte(data.ColorDeviceB.UV),
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

	success, palette := fixture.Visuals.GetPalettes().GetById(data.PaletteId)
	if !success {
		http.Error(w, "Palette not found.", http.StatusBadRequest)
		return
	}

	durationCycles := int(float64(data.DurationMs) * fadeDurationMsToCycles)
	if durationCycles <= 0 || durationCycles > 0xFFFF {
		http.Error(w, "Fade duration outside of range.", http.StatusBadRequest)
		return
	}

	if data.PaletteShift < 0 || data.PaletteShift > 7 {
		http.Error(w, "Palette shift outside of range.", http.StatusBadRequest)
		return
	}

	block, err := BoxiBus.CreateLightingPaletteFade(palette.Colors, uint16(durationCycles), byte(data.PaletteShift), false)
	if err != nil {
		http.Error(w, fmt.Sprintf("Instruction couldn't be created. %s", err), http.StatusInternalServerError)
		return
	}

	instruction := Lightshow.LightingInstruction{
		MessageBlock: block,
	}

	fixture.Visuals.SetLightingOverwrite(&instruction)
}

func (fixture Fixture) HandleSetLightingOverridePaletteSwitchApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var data lightingInstructionPaletteSwitch

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid parameters. %s", err), http.StatusBadRequest)
		return
	}

	success, palette := fixture.Visuals.GetPalettes().GetById(data.PaletteId)
	if !success {
		http.Error(w, "Palette not found.", http.StatusBadRequest)
		return
	}

	if data.PaletteShift < 0 || data.PaletteShift > 7 {
		http.Error(w, "Palette shift outside of range.", http.StatusBadRequest)
		return
	}

	block, err := BoxiBus.CreateLightingPaletteSwitch(palette.Colors, byte(data.PaletteShift), false)
	if err != nil {
		http.Error(w, fmt.Sprintf("Instruction couldn't be created. %s", err), http.StatusInternalServerError)
		return
	}

	instruction := Lightshow.LightingInstruction{
		MessageBlock: block,
	}

	fixture.Visuals.SetLightingOverwrite(&instruction)
}

func (fixture Fixture) HandleSetLightingOverridePaletteBrightnessFlashApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var data lightingInstructionPaletteBrightnessFlash

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid parameters. %s", err), http.StatusBadRequest)
		return
	}

	success, palette := fixture.Visuals.GetPalettes().GetById(data.PaletteId)
	if !success {
		http.Error(w, "Palette not found.", http.StatusBadRequest)
		return
	}

	fadeOutCycles := int(float64(data.DurationMs) * fadeDurationMsToCycles)
	if fadeOutCycles <= 0 || fadeOutCycles > 0xFFFF {
		http.Error(w, "Fade out duration outside of range.", http.StatusBadRequest)
		return
	}

	if data.TargetBrightness <= 0 || data.TargetBrightness > 0xFFFF {
		http.Error(w, "Target brightness duration outside of range.", http.StatusBadRequest)
		return
	}

	if data.PaletteShift < 0 || data.PaletteShift > 7 {
		http.Error(w, "Palette shift outside of range.", http.StatusBadRequest)
		return
	}

	block, err := BoxiBus.CreateLightingPaletteBrightnessFlash(palette.Colors, uint16(fadeOutCycles), byte(data.PaletteShift), false)
	if err != nil {
		http.Error(w, fmt.Sprintf("Instruction couldn't be created. %s", err), http.StatusInternalServerError)
		return
	}

	instruction := Lightshow.LightingInstruction{
		MessageBlock: block,
	}

	fixture.Visuals.SetLightingOverwrite(&instruction)
}

func (fixture Fixture) HandleSetLightingOverridePaletteBrightnessFlashApi(w http.ResponseWriter, r *http.Request) {
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

	success, palette := fixture.Visuals.GetPalettes().GetById(data.PaletteId)
	if !success {
		http.Error(w, "Palette not found.", http.StatusBadRequest)
		return
	}

	fadeOutCycles := int(float64(data.DurationMs) * fadeDurationMsToCycles)
	if fadeOutCycles <= 0 || fadeOutCycles > 0xFFFF {
		http.Error(w, "Fade out duration outside of range.", http.StatusBadRequest)
		return
	}

	if data.PaletteShift < 0 || data.PaletteShift > 7 {
		http.Error(w, "Palette shift outside of range.", http.StatusBadRequest)
		return
	}

	block, err := BoxiBus.CreateLightingPaletteHueFlash(palette.Colors, uint16(fadeOutCycles), byte(data.PaletteShift), false)
	if err != nil {
		http.Error(w, fmt.Sprintf("Instruction couldn't be created. %s", err), http.StatusInternalServerError)
		return
	}

	instruction := Lightshow.LightingInstruction{
		MessageBlock: block,
	}

	fixture.Visuals.SetLightingOverwrite(&instruction)
}
