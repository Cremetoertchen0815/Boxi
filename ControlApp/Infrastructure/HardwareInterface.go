package Infrastructure

import (
	"ControlApp/BoxiBus"
	"ControlApp/Display"
)

type HardwareInterface interface {
	GetConnectedDisplays() []Display.ServerDisplay
	SendLightingInstruction(block BoxiBus.MessageBlock)
	SendAnimationInstruction(animation Display.AnimationId, displays []Display.ServerDisplay)
	SendTextInstruction(text string, displays []Display.ServerDisplay)
	SendBrightnessChange(brightness *float64, blinkSpeed uint16)
	SendBeatToDisplay(force bool)
	UploadAnimation(id Display.AnimationId)
}
