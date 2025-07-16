package Lightshow

import (
	"ControlApp/BoxiBus"
	"encoding/json"
	"log"
	"os"
	"sort"
	"sync"
)

type PaletteManager struct {
	palettes   map[uint32]Palette
	accessLock *sync.Mutex
}

type Palette struct {
	Id     uint32
	Name   string
	Colors []BoxiBus.Color
	Moods  []LightingMood
}

const palettesConfigPath = "Configuration/palettes.json"
const palettesConfigBackupPath = "Configuration/palettes_backup.json"

func LoadPalettes() *PaletteManager {
	config, err := loadConfiguration[map[uint32]Palette](palettesConfigPath)

	if err != nil {
		config, err = loadConfiguration[map[uint32]Palette](palettesConfigBackupPath)
	}

	if err != nil {
		log.Fatalf("Config file for palettes could not be accessed! %s", err)
	}

	return &PaletteManager{
		palettes:   config,
		accessLock: &sync.Mutex{},
	}
}

func (manager *PaletteManager) storeConfiguration() {
	_ = os.Remove(palettesConfigPath)

	configFile, err := os.OpenFile(palettesConfigPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)

	if err != nil {
		log.Fatalf("Config file for palettes could not be opened for writing! %s", err)
	}

	defer func(configFile *os.File) {
		_ = configFile.Close()
	}(configFile)

	jsonParser := json.NewEncoder(configFile)
	err = jsonParser.Encode(manager.palettes)
	if err != nil {
		log.Fatalf("Configuration for palettes could be JSON encoded! %s", err)
	}
}

func (manager *PaletteManager) GetPalettesForMood(mood LightingMood) []Palette {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

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
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	for _, palette := range manager.palettes {
		if palette.Id == id {
			return true, palette
		}
	}

	return false, Palette{}
}

func (manager *PaletteManager) GetAll() []Palette {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	var palettes []Palette
	for _, palette := range manager.palettes {
		palettes = append(palettes, palette)
	}

	sort.Slice(palettes, func(i, j int) bool {
		return palettes[i].Id < palettes[j].Id
	})

	return palettes
}

func (manager *PaletteManager) SetPalette(palette Palette) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	manager.palettes[palette.Id] = palette
	manager.storeConfiguration()
}

func (manager *PaletteManager) RemovePalette(paletteId uint32) {
	manager.accessLock.Lock()
	defer manager.accessLock.Unlock()

	delete(manager.palettes, paletteId)
	manager.storeConfiguration()
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
