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
	Id     uint32
	Name   string
	Colors []BoxiBus.Color
	Moods  []LightingMood
}

func LoadPalettes() *PaletteManager {
	return &PaletteManager{
		palettes: []Palette{
			{
				Id:   0,
				Name: "Rainbow",
				Colors: []BoxiBus.Color{
					{255, 0, 0, 0, 0, 0},
					{255, 255, 0, 0, 0, 0},
					{0, 255, 0, 0, 0, 0},
					{0, 255, 255, 0, 0, 0},
					{0, 0, 255, 0, 0, 0},
					{255, 0, 255, 0, 0, 0},
				},
				Moods: []LightingMood{Happy, Regular, Party},
			},
			{
				Id:   1,
				Name: "Cyberpunk",
				Colors: []BoxiBus.Color{
					{0, 255, 153, 0, 0, 0},
					{77, 156, 200, 24, 0, 0},
					{0, 30, 255, 0, 0, 128},
					{150, 0, 200, 0, 0, 128},
					{128, 0, 128, 10, 0, 255},
				},
				Moods: []LightingMood{Moody, Regular},
			},
			{
				Id:   2,
				Name: "Pleasant",
				Colors: []BoxiBus.Color{
					{108, 200, 25, 0, 0, 0},
					{255, 255, 25, 0, 50, 0},
					{255, 130, 50, 0, 255, 32},
					{255, 64, 80, 0, 0, 32},
					{100, 100, 255, 0, 0, 255},
					{52, 128, 255, 40, 0, 64},
				},
				Moods: []LightingMood{Happy, Regular},
			},
			{
				Id:   3,
				Name: "Retro",
				Colors: []BoxiBus.Color{
					{0, 0, 0, 255, 128, 0},
					{180, 180, 0, 0, 255, 0},
					{180, 20, 80, 0, 0, 128},
					{20, 0, 64, 0, 255, 255},
					{0, 255, 0, 0, 255, 0},
					{255, 0, 0, 0, 0, 255},
				},
				Moods: []LightingMood{Regular, Party},
			},
		},
		accessLock: &sync.Mutex{},
	}
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

func (manager *PaletteManager) GetById(id uint32) (bool, Palette) {
	for _, palette := range manager.palettes {
		if palette.Id == id {
			return true, palette
		}
	}

	return false, Palette{}
}

func getDefaultPalettes() []Palette {
	return []Palette{
		{0, "Default Rainbow",
			[]BoxiBus.Color{
				{255, 0, 0, 0, 0, 0},
				{255, 255, 0, 0, 0, 0},
				{0, 255, 0, 0, 0, 0},
				{0, 255, 255, 0, 0, 0},
				{0, 0, 255, 0, 0, 0},
				{255, 0, 255, 0, 0, 0},
			}, []LightingMood{Moody, Happy, Regular, Party}},
	}
}
