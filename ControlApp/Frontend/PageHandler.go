package Frontend

import (
	"ControlApp/Api"
	"ControlApp/BoxiBus"
	"ControlApp/Lightshow"
	"fmt"
	"net/http"
)

type startPageInformation struct {
	ScaffoldInformation
	Mood              int
	Nsfw              bool
	Brightness        int
	ConnectedDisplays string
}

type overridePageInformation struct {
	ScaffoldInformation
	LightingOverride        bool
	LightingMode            int
	LightingShowColorA      bool
	LightingShowColorB      bool
	LightingColorA          string
	LightingColorB          string
	LightingShowPalettes    bool
	LightingPalettes        []Lightshow.Palette
	LightingPaletteId       uint32
	LightingShowDuration    bool
	LightingDurationValue   int
	LightingShowBrightness  bool
	LightingBrightnessValue int
	LightingShowFrequency   bool
	LightingFrequencyValue  int
	LightingShowSpeed       bool
	LightingSpeedValue      int
	LightingShowShift       bool
	LightingShiftValue      int
}

func (Me PageProvider) HandleStartPage(w http.ResponseWriter, r *http.Request) {
	//Fetch scaffold data from context
	scaffoldData := GetScaffoldData(r)

	//Create data structure
	mood := int(Me.Data.Visuals.GetConfiguration().Mood)
	isNsfw := Me.Data.Visuals.GetConfiguration().AllowNsfw
	brightness := int(Me.Data.Visuals.GetBrightness() * 100)
	displays := fmt.Sprintf("%+v", Me.Data.Hardware.GetConnectedDisplays())
	startData := startPageInformation{scaffoldData, mood, isNsfw, brightness, displays}

	//Disable caching
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	//Execute template
	err := Me.startPage.Execute(w, startData)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (Me PageProvider) HandleOverridesPage(w http.ResponseWriter, r *http.Request) {
	//Fetch scaffold data from context
	scaffoldData := GetScaffoldData(r)

	mode := BoxiBus.LightingModeId(Me.Data.OverrideLightingCurrent.Mode)
	showColorA := mode == BoxiBus.SetColor || mode == BoxiBus.FadeToColor || mode == BoxiBus.Strobe
	showColorB := mode == BoxiBus.SetColor || mode == BoxiBus.FadeToColor
	showPalette := mode == BoxiBus.PaletteFade || mode == BoxiBus.PaletteSwitch || mode == BoxiBus.PaletteBrightnessFlash || mode == BoxiBus.PaletteHueFlash
	showDuration := mode == BoxiBus.FadeToColor || mode == BoxiBus.PaletteFade
	showSpeed := mode == BoxiBus.PaletteBrightnessFlash || mode == BoxiBus.PaletteHueFlash
	showShift := mode == BoxiBus.PaletteFade || mode == BoxiBus.PaletteSwitch || mode == BoxiBus.PaletteBrightnessFlash || mode == BoxiBus.PaletteHueFlash
	showFrequency := mode == BoxiBus.Strobe

	data := overridePageInformation{
		ScaffoldInformation:     scaffoldData,
		LightingOverride:        Me.Data.OverrideLightingCurrent.Enable,
		LightingMode:            Me.Data.OverrideLightingCurrent.Mode,
		LightingShowColorA:      showColorA,
		LightingColorA:          getColorString(Me.Data.OverrideLightingCurrent.ColorDeviceA),
		LightingShowColorB:      showColorB,
		LightingColorB:          getColorString(Me.Data.OverrideLightingCurrent.ColorDeviceB),
		LightingShowPalettes:    showPalette,
		LightingPalettes:        Me.Data.Visuals.GetPalettes().GetAll(),
		LightingPaletteId:       Me.Data.OverrideLightingCurrent.PaletteId,
		LightingShowDuration:    showSpeed,
		LightingDurationValue:   Me.Data.OverrideLightingCurrent.DurationMs,
		LightingShowBrightness:  showDuration,
		LightingBrightnessValue: Me.Data.OverrideLightingCurrent.TargetBrightness,
		LightingShowFrequency:   showFrequency,
		LightingFrequencyValue:  Me.Data.OverrideLightingCurrent.FrequencyHz,
		LightingShowShift:       showShift,
		LightingShiftValue:      Me.Data.OverrideLightingCurrent.PaletteShift,
		LightingShowSpeed:       showSpeed,
		LightingSpeedValue:      Me.Data.OverrideLightingCurrent.Speed,
	}

	//Disable caching
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	//Execute template
	err := Me.overridesPage.Execute(w, data)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func getColorString(color Api.Color) string {
	return fmt.Sprintf("%d,%d,%d,%d,%d,%d", color.R, color.G, color.B, color.W, color.A, color.UV)
}

func (Me PageProvider) HandleAnimationPage(w http.ResponseWriter, r *http.Request) {
	//Fetch scaffold data from context
	scaffoldData := GetScaffoldData(r)

	//Disable caching
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	//Execute template
	err := Me.animationsPage.Execute(w, scaffoldData)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (Me PageProvider) HandlePalettesPage(w http.ResponseWriter, r *http.Request) {
	//Fetch scaffold data from context
	scaffoldData := GetScaffoldData(r)

	//Disable caching
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	//Execute template
	err := Me.palettesPage.Execute(w, scaffoldData)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (Me PageProvider) HandleAutoPage(w http.ResponseWriter, r *http.Request) {
	//Fetch scaffold data from context
	scaffoldData := GetScaffoldData(r)

	//Disable caching
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	//Execute template
	err := Me.autoPage.Execute(w, scaffoldData)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}
