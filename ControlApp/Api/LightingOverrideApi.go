package Api

import (
	"ControlApp/BoxiBus"
	"ControlApp/Lightshow"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	fadeDurationMsToCycles    float64 = 0.130
	strobeFrequencyMultiplier float64 = 24
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
	TargetBrightness int    `json:"targetBrightness"`
	Speed            int    `json:"speed"`
	PaletteShift     int    `json:"paletteShift"`
}

type lightingInstructionPaletteHueFlash struct {
	PaletteId    uint32 `json:"paletteId"`
	Speed        int    `json:"speed"`
	PaletteShift int    `json:"paletteShift"`
}

type lightingInstructionStrobe struct {
	Color       color   `json:"color"`
	FrequencyHz float64 `json:"frequency"`
}

func (fixture Fixture) HandleSetLightingOverrideAutoApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	fixture.Visuals.SetLightingOverwrite(nil)
}

func (fixture Fixture) HandleSetLightingOverrideOffApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	instruction := Lightshow.LightingInstruction{
		MessageBlock: BoxiBus.CreateLightingOff(false),
	}

	fixture.Visuals.SetLightingOverwrite(&instruction)
}

func (fixture Fixture) HandleSetLightingOverrideSetColorApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
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
	if r.Method != "POST" {
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
	if r.Method != "POST" {
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
	if r.Method != "POST" {
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
	if r.Method != "POST" {
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

	if data.Speed <= 0 || data.Speed > 0xFFFF {
		http.Error(w, "Fade out duration outside of range.", http.StatusBadRequest)
		return
	}

	if data.TargetBrightness < 0 || data.TargetBrightness > 0xFF {
		http.Error(w, "Target brightness duration outside of range.", http.StatusBadRequest)
		return
	}

	if data.PaletteShift < 0 || data.PaletteShift > 7 {
		http.Error(w, "Palette shift outside of range.", http.StatusBadRequest)
		return
	}

	block, err := BoxiBus.CreateLightingPaletteBrightnessFlash(palette.Colors, uint16(data.Speed), byte(data.TargetBrightness), byte(data.PaletteShift), false)
	if err != nil {
		http.Error(w, fmt.Sprintf("Instruction couldn't be created. %s", err), http.StatusInternalServerError)
		return
	}

	instruction := Lightshow.LightingInstruction{
		MessageBlock: block,
	}

	fixture.Visuals.SetLightingOverwrite(&instruction)
}

func (fixture Fixture) HandleSetLightingOverridePaletteHueFlashApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var data lightingInstructionPaletteHueFlash

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

	if data.Speed <= 0 || data.Speed > 0xFF {
		http.Error(w, "Fade out duration outside of range.", http.StatusBadRequest)
		return
	}

	if data.PaletteShift < 0 || data.PaletteShift > 7 {
		http.Error(w, "Palette shift outside of range.", http.StatusBadRequest)
		return
	}

	block, err := BoxiBus.CreateLightingPaletteHueFlash(palette.Colors, uint16(data.Speed), byte(data.PaletteShift), false)
	if err != nil {
		http.Error(w, fmt.Sprintf("Instruction couldn't be created. %s", err), http.StatusInternalServerError)
		return
	}

	instruction := Lightshow.LightingInstruction{
		MessageBlock: block,
	}

	fixture.Visuals.SetLightingOverwrite(&instruction)
}

func (fixture Fixture) HandleSetLightingOverrideStrobeApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var data lightingInstructionStrobe

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid parameters. %s", err), http.StatusBadRequest)
		return
	}

	frequency := int((1 / float64(data.FrequencyHz)) * strobeFrequencyMultiplier)
	if frequency <= 0 || frequency > 0xF {
		http.Error(w, "frequency outside of range.", http.StatusBadRequest)
		return
	}

	if !isColorValid(data.Color) {
		http.Error(w, "Color is invalid.", http.StatusBadRequest)
		return
	}

	color := BoxiBus.Color{
		Red:         byte(data.Color.R),
		Green:       byte(data.Color.G),
		Blue:        byte(data.Color.B),
		White:       byte(data.Color.W),
		Amber:       byte(data.Color.A),
		UltraViolet: byte(data.Color.UV),
	}

	instruction := Lightshow.LightingInstruction{
		MessageBlock: BoxiBus.CreateLightingStrobe(color, uint16(frequency), 0, false),
	}

	fixture.Visuals.SetLightingOverwrite(&instruction)
}
