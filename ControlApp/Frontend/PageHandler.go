package Frontend

import (
	"ControlApp/Api"
	"ControlApp/BoxiBus"
	"ControlApp/Display"
	"ControlApp/Lightshow"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type startPageInformation struct {
	ScaffoldInformation
	Mood              int
	Nsfw              bool
	Brightness        int
	ConnectedDisplays string
}

type animationInformation struct {
	Api.ScreenOverrideAnimationInstance
	ScreenNumber int
}

type textInformation struct {
	Api.ScreenTextInstance
	ScreenNumber int
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
	AnimationsOverride      bool
	Animations              []Lightshow.Animation
	AnimationsSelected      []animationInformation
	AnimationsFadeout       int
	TextOverride            bool
	TextOverrideValues      []textInformation
}

type palettePageInformation struct {
	ScaffoldInformation
	Palettes []Lightshow.Palette
}

type autoModePageInformation struct {
	ScaffoldInformation
	Api.AutoModeConfig
}

type animationsPageInformation struct {
	ScaffoldInformation
	Animations []animationInstance
}

type animationInstance struct {
	Id        uint32
	Name      string
	Details   string
	Thumbnail string
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
	showBrightness := mode == BoxiBus.PaletteBrightnessFlash
	showSpeed := mode == BoxiBus.PaletteBrightnessFlash || mode == BoxiBus.PaletteHueFlash
	showShift := mode == BoxiBus.PaletteFade || mode == BoxiBus.PaletteSwitch || mode == BoxiBus.PaletteBrightnessFlash || mode == BoxiBus.PaletteHueFlash
	showFrequency := mode == BoxiBus.Strobe

	var animations []animationInformation
	for _, anim := range Me.Data.OverrideAnimationCurrent.Animations {
		number := 1
		switch anim.ScreenIndex {
		case Display.Boxi1D2:
			number = 2
		case Display.Boxi2D1:
			number = 3
		case Display.Boxi2D2:
			number = 4
		}
		infoStruct := animationInformation{anim, number}
		animations = append(animations, infoStruct)
	}

	allAnimations := Me.Data.Visuals.GetAnimations().GetAll()
	sort.Slice(allAnimations, func(i, j int) bool {
		return strings.ToLower(allAnimations[i].Name) < strings.ToLower(allAnimations[j].Name)
	})

	anyTextOverwrites := false
	var texts []textInformation
	for _, text := range Me.Data.OverrideTextsCurrent.Texts {
		textContent := text.Text
		//Check if the text overwrite is empty
		if strings.TrimSpace(textContent) == "" {
			textContent = ""
		} else {
			anyTextOverwrites = true
		}

		number := 1
		switch text.ScreenIndex {
		case Display.Boxi1D2:
			number = 2
		case Display.Boxi2D1:
			number = 3
		case Display.Boxi2D2:
			number = 4
		}

		coreData := Api.ScreenTextInstance{ScreenIndex: text.ScreenIndex, Text: textContent}
		texts = append(texts, textInformation{coreData, number})
	}

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
		LightingShowDuration:    showDuration,
		LightingDurationValue:   Me.Data.OverrideLightingCurrent.DurationMs,
		LightingShowBrightness:  showBrightness,
		LightingBrightnessValue: Me.Data.OverrideLightingCurrent.TargetBrightness,
		LightingShowFrequency:   showFrequency,
		LightingFrequencyValue:  Me.Data.OverrideLightingCurrent.FrequencyHz,
		LightingShowShift:       showShift,
		LightingShiftValue:      Me.Data.OverrideLightingCurrent.PaletteShift,
		LightingShowSpeed:       showSpeed,
		LightingSpeedValue:      Me.Data.OverrideLightingCurrent.Speed,
		AnimationsOverride:      !Me.Data.OverrideAnimationCurrent.ResetScreens,
		Animations:              allAnimations,
		AnimationsSelected:      animations,
		AnimationsFadeout:       Me.Data.OverrideAnimationCurrent.FadeoutSpeed,
		TextOverride:            anyTextOverwrites,
		TextOverrideValues:      texts,
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

	var animations []animationInstance
	for _, animation := range Me.Data.Visuals.GetAnimations().GetAll() {
		moodStr := "Unknown"
		switch animation.Mood {
		case Lightshow.Happy:
			moodStr = "Happy"
			break
		case Lightshow.Moody:
			moodStr = "Moody"
			break
		case Lightshow.Regular:
			moodStr = "Regular"
			break
		case Lightshow.Party:
			moodStr = "Party"
			break
		}

		nsfwStr := "Not NSFW"
		if animation.IsNsfw {
			nsfwStr = "NSFW"
		}

		aniInstance := animationInstance{
			Id:        uint32(animation.Id),
			Name:      animation.Name,
			Details:   moodStr + ", " + nsfwStr,
			Thumbnail: fmt.Sprintf("/static/thumbs/%d.png", animation.Id),
		}
		animations = append(animations, aniInstance)
	}

	//Fetch scaffold data from context
	scaffoldData := GetScaffoldData(r)
	templateData := animationsPageInformation{scaffoldData, animations}

	//Disable caching
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	//Execute template
	err := Me.animationsPage.Execute(w, templateData)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (Me PageProvider) HandlePalettesPage(w http.ResponseWriter, r *http.Request) {
	//Fetch scaffold data from context
	scaffoldData := GetScaffoldData(r)
	templateData := palettePageInformation{scaffoldData, Me.Data.Visuals.GetPalettes().GetAll()}

	//Disable caching
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	//Execute template
	err := Me.palettesPage.Execute(w, templateData)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (Me PageProvider) HandleAutoPage(w http.ResponseWriter, r *http.Request) {
	//Fetch scaffold data from context
	scaffoldData := GetScaffoldData(r)
	configData := Api.AutoModeConfig{
		StrobeChance:               0,
		HueShiftChance:             0,
		FadeToColorCycles:          0,
		PaletteFadeCycles:          0,
		FlashFadeoutSpeed:          0,
		HueFlashFadeoutSpeed:       0,
		StrobeFrequency:            0,
		FlashTargetBrightness:      0,
		FlashHueShift:              0,
		MinTimeBetweenBeatsSec:     0,
		LightingCalmModeBoringSec:  0,
		AnimationCalmModeBoringSec: 0,
		CalmLightingTiming:         Api.TimingConstraint{},
		RhythmicLightingTiming:     Api.TimingConstraint{},
		FranticLightingTiming:      Api.TimingConstraint{},
		CalmAnimationsTiming:       Api.TimingConstraint{},
		RhythmicAnimationsTiming:   Api.TimingConstraint{},
		FranticAnimationsTiming:    Api.TimingConstraint{},
	}

	//configuration := fixture.Data.Visuals.GetConfiguration()
	//	configuration.StrobeChance = data.StrobeChance
	//	configuration.HueShiftChance = data.HueShiftChance
	//	configuration.FadeToColorCycles = data.FadeToColorCycles
	//	configuration.PaletteFadeCycles = data.PaletteFadeCycles
	//	configuration.FlashFadeoutSpeed = data.FlashFadeoutSpeed
	//	configuration.HueFlashFadeoutSpeed = data.HueFlashFadeoutSpeed
	//	configuration.StrobeFrequency = data.StrobeFrequency
	//	configuration.FlashTargetBrightness = data.FlashTargetBrightness
	//	configuration.FlashHueShift = data.FlashHueShift
	//	configuration.MinTimeBetweenBeats = time.Duration(float64(time.Second) * data.MinTimeBetweenBeatsSec)
	//	configuration.LightingCalmModeBoring = time.Duration(float64(time.Second) * data.LightingCalmModeBoringSec)
	//	configuration.AnimationCalmModeBoring = time.Duration(float64(time.Second) * data.AnimationCalmModeBoringSec)
	//	configuration.LightingModeTiming[Lightshow.Calm] = getConstraint(data.CalmLightingTiming)
	//	configuration.LightingModeTiming[Lightshow.Rhythmic] = getConstraint(data.RhythmicLightingTiming)
	//	configuration.LightingModeTiming[Lightshow.Frantic] = getConstraint(data.FranticLightingTiming)
	//	configuration.AnimationModeTiming[Lightshow.Calm] = getConstraint(data.CalmAnimationsTiming)
	//	configuration.AnimationModeTiming[Lightshow.Rhythmic] = getConstraint(data.RhythmicAnimationsTiming)
	//	configuration.AnimationModeTiming[Lightshow.Frantic] = getConstraint(data.FranticAnimationsTiming)

	templateData := autoModePageInformation{
		ScaffoldInformation: scaffoldData,
		AutoModeConfig:      configData,
	}

	//Disable caching
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	//Execute template
	err := Me.autoPage.Execute(w, templateData)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}
