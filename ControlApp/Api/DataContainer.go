package Api

import (
	"ControlApp/Display"
	"ControlApp/Infrastructure"
	"ControlApp/Lightshow"
)

type DataContainer struct {
	Hardware                 Infrastructure.HardwareInterface
	Visuals                  *Lightshow.VisualManager
	OverrideLightingCurrent  LightingInstructionTotal
	OverrideAnimationCurrent ScreenOverrideAnimationProperties
	OverrideTextsCurrent     ScreenOverrideTextProperties
}

func CreateDataContainer(hardware Infrastructure.HardwareInterface, visuals *Lightshow.VisualManager) *DataContainer {
	result := DataContainer{
		hardware, visuals, LightingInstructionTotal{
			Enable:           false,
			ApplyOnBeat:      false,
			Mode:             0,
			ColorDeviceA:     Color{R: 255},
			ColorDeviceB:     Color{G: 255},
			PaletteId:        0,
			DurationMs:       2000,
			PaletteShift:     0,
			Speed:            40,
			TargetBrightness: 15,
			FrequencyHz:      12,
		},
		ScreenOverrideAnimationProperties{
			Animations: []ScreenOverrideAnimationInstance{
				{ScreenIndex: Display.Boxi1D1, AnimationId: 0},
				{ScreenIndex: Display.Boxi1D2, AnimationId: 0},
				{ScreenIndex: Display.Boxi2D1, AnimationId: 0},
				{ScreenIndex: Display.Boxi2D2, AnimationId: 0},
			},
			FadeoutSpeed: 0,
			ResetScreens: true,
		},
		ScreenOverrideTextProperties{
			Texts: []ScreenTextInstance{
				{ScreenIndex: Display.Boxi1D1, Text: " "},
				{ScreenIndex: Display.Boxi1D2, Text: " "},
				{ScreenIndex: Display.Boxi2D1, Text: " "},
				{ScreenIndex: Display.Boxi2D2, Text: " "},
			},
		},
	}

	return &result
}
