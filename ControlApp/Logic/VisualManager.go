package Logic

import (
	"ControlApp/BoxiBus"
	"ControlApp/Lightshow"
)

type VisualManager struct {
	Running     bool
	autoContext Lightshow.AutoModeContext
}

type lightingInstructions struct {
	Led *BoxiBus.MessageBlock //The message block that will be immediately sent to
}

func CreateLightingManager(hardwareManager HardwareManager) *VisualManager {
	autoContext := Lightshow.CreateAutoMode(hardwareManager, loadConfiguration())

	return result
}

func loadConfiguration() Lightshow.AutoModeConfiguration {
	return Lightshow.AutoModeConfiguration{}
}
