package Lightshow

import "ControlApp/BoxiBus"

type ModeCharacter uint8

const (
	Calm ModeCharacter = iota
	Rhythmic
	Frantic
	Unknown
)

func getLightingModeCharacter(modeId BoxiBus.LightingModeId) ModeCharacter {
	switch modeId {
	case BoxiBus.SetColor, BoxiBus.FadeToColor, BoxiBus.PaletteFade:
		return Calm
	case BoxiBus.PaletteSwitch, BoxiBus.PaletteBrightnessFlash, BoxiBus.PaletteHueFlash:
		return Rhythmic
	case BoxiBus.Strobe:
		return Frantic
	}

	return Unknown
}
