package Logic

import (
	"github.com/stianeikeland/go-rpio/v4"
	"time"
)

type LightingManager struct {
	Configuration LightingConfig
	hardware      HardwareManager
	Running       bool
	inputPin      rpio.Pin
	lastBeat      *time.Time
}

type LightingConfig struct {
	MinTimeBetweenBeats time.Duration
}

const soundInputPin = 16
const loopDelayMs = 5

func CreateLightingManager(hardwareManager HardwareManager) *LightingManager {
	pin := rpio.Pin(soundInputPin)
	pin.Input()

	result := &LightingManager{loadConfiguration(), hardwareManager, true, pin, nil}

	go result.workLighting()
	return result
}

func loadConfiguration() LightingConfig {
	return LightingConfig{}
}

func (manager *LightingManager) workLighting() {
	for manager.Running {
		time.Sleep(loopDelayMs * time.Millisecond)

		isBeat := manager.inputPin.Read() == rpio.High
		now := time.Now()

		if !isBeat || (manager.lastBeat != nil && now.Before(manager.lastBeat.Add(manager.Configuration.MinTimeBetweenBeats))) {
			continue
		}

		manager.lastBeat = &now
	}
}
