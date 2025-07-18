package Api

import (
	"ControlApp/BoxiBus"
	"ControlApp/Infrastructure"
	"ControlApp/Lightshow"
	"encoding/json"
	"fmt"
	"net/http"
)

type LightingInstructionTotal struct {
	Enable           bool   `json:"enable"`
	ApplyOnBeat      bool   `json:"onBeat"`
	Mode             int    `json:"mode"`
	ColorDeviceA     Color  `json:"colorA"`
	ColorDeviceB     Color  `json:"colorB"`
	PaletteId        uint32 `json:"paletteId"`
	DurationMs       int    `json:"duration"`
	PaletteShift     int    `json:"paletteShift"`
	Speed            int    `json:"speed"`
	TargetBrightness int    `json:"targetBrightness"`
	FrequencyHz      int    `json:"frequency"`
}

func (fixture Fixture) HandleSetLightingOverrideAutoApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var data LightingInstructionTotal

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid parameters. %s", err), http.StatusBadRequest)
		return
	}

	fixture.Data.OverrideLightingCurrent = data

	if !data.Enable {
		fixture.Data.Visuals.SetLightingOverwrite(nil)
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

	durationCycles := int(float64(data.DurationMs) * Infrastructure.FadeDurationMsToCycles)
	if durationCycles <= 0 || durationCycles > 0xFFFF {
		http.Error(w, "Fade duration outside of range.", http.StatusBadRequest)
		return
	}

	success, palette := fixture.Data.Visuals.GetPalettes().GetById(data.PaletteId)
	if !success {
		http.Error(w, "Palette not found.", http.StatusBadRequest)
		return
	}

	if data.PaletteShift < 0 || data.PaletteShift > 7 {
		http.Error(w, "Palette shift outside of range.", http.StatusBadRequest)
		return
	}

	frequency := int((1 / float64(data.FrequencyHz)) * Infrastructure.StrobeFrequencyMultiplier)
	if frequency <= 0 || frequency > 0xF {
		http.Error(w, "frequency outside of range.", http.StatusBadRequest)
		return
	}

	var instruction Lightshow.LightingInstruction

	switch data.Mode {
	case 0:
		instruction = Lightshow.LightingInstruction{
			MessageBlock: BoxiBus.CreateLightingOff(data.ApplyOnBeat),
		}
		break
	case 1:
		instruction = Lightshow.LightingInstruction{
			MessageBlock: BoxiBus.CreateLightingSetColor(color1, color2, data.ApplyOnBeat),
		}
		break
	case 2:
		instruction = Lightshow.LightingInstruction{
			MessageBlock: BoxiBus.CreateLightingFadeToColor(color1, color2, uint16(durationCycles), data.ApplyOnBeat),
		}
		break
	case 3:
		block, err := BoxiBus.CreateLightingPaletteFade(palette.Colors, uint16(durationCycles), byte(data.PaletteShift), data.ApplyOnBeat)
		if err != nil {
			http.Error(w, fmt.Sprintf("Instruction couldn't be created. %s", err), http.StatusInternalServerError)
			return
		}

		instruction = Lightshow.LightingInstruction{
			MessageBlock: block,
		}
		break
	case 4:
		block, err := BoxiBus.CreateLightingPaletteSwitch(palette.Colors, byte(data.PaletteShift), data.ApplyOnBeat)
		if err != nil {
			http.Error(w, fmt.Sprintf("Instruction couldn't be created. %s", err), http.StatusInternalServerError)
			return
		}

		instruction = Lightshow.LightingInstruction{
			MessageBlock: block,
		}
		break
	case 5:
		block, err := BoxiBus.CreateLightingPaletteBrightnessFlash(palette.Colors, uint16(data.Speed), byte(data.TargetBrightness), byte(data.PaletteShift), data.ApplyOnBeat)
		if err != nil {
			http.Error(w, fmt.Sprintf("Instruction couldn't be created. %s", err), http.StatusInternalServerError)
			return
		}

		instruction = Lightshow.LightingInstruction{
			MessageBlock: block,
		}
		break
	case 6:

		block, err := BoxiBus.CreateLightingPaletteHueFlash(palette.Colors, uint16(data.Speed), byte(data.PaletteShift), data.ApplyOnBeat)
		if err != nil {
			http.Error(w, fmt.Sprintf("Instruction couldn't be created. %s", err), http.StatusInternalServerError)
			return
		}

		instruction = Lightshow.LightingInstruction{
			MessageBlock: block,
		}
		break
	case 7:
		instruction = Lightshow.LightingInstruction{
			MessageBlock: BoxiBus.CreateLightingStrobe(color1, uint16(frequency), 0, data.ApplyOnBeat),
		}
		break
	}

	fixture.Data.Visuals.SetLightingOverwrite(&instruction)
}
