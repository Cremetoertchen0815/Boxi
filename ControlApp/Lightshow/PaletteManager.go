package Lightshow

import (
	"ControlApp/BoxiBus"
	"sync"
)

type PaletteManager struct {
	palettes   []Palette
	accessLock *sync.Mutex
}

type Palette struct {
	Colors []BoxiBus.Color
	Moods  []LightingMood
}

func LoadPalettes() *PaletteManager {
	return &PaletteManager{}
}

func (manager *PaletteManager) GetPalettesForMood(mood LightingMood) []Palette {
	var result []Palette
	for _, palette := range manager.palettes {
		for _, paletteMode := range palette.Moods {
			if paletteMode != mood {
				continue
			}

			result = append(result, palette)
			break
		}
	}

	return result
}

func getDefaultPalettes() []Palette {
	return []Palette{
		{[]BoxiBus.Color{
			{255, 0, 0, 0, 0, 0},
			{255, 255, 0, 0, 0, 0},
			{0, 255, 0, 0, 0, 0},
			{0, 255, 255, 0, 0, 0},
			{0, 0, 255, 0, 0, 0},
			{255, 0, 255, 0, 0, 0},
		}, []LightingMood{Moody, Happy, Regular, Party}},
	}
}
