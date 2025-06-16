package Infrastructure

import (
	"ControlApp/BoxiBus"
	"ControlApp/Display"
	"log"
)

type DebugHardwareManager int

func (manager DebugHardwareManager) GetConnectedDisplays() []Display.ServerDisplay {
	return []Display.ServerDisplay{Display.ServerDisplay(1)}
}

func (manager DebugHardwareManager) SendLightingInstruction(block BoxiBus.MessageBlock) {
	log.Printf("Lighting instruction sent: %+v \n", block)
}

func (manager DebugHardwareManager) SendAnimationInstruction(animation Display.AnimationId, displays []Display.ServerDisplay) {
	log.Printf("Animation instruction sent, animation: %d, displays: %+v \n", animation, displays)
}

func (manager DebugHardwareManager) SendTextInstruction(text string, displays []Display.ServerDisplay) {
	log.Printf("Animation instruction sent, text: %s, displays: %+v \n", text, displays)
}

func (manager DebugHardwareManager) SendBrightnessChange(brightness *float64, blinkSpeed uint16) {
	log.Printf("Animation instruction sent, brightness: %+v, speed: %d \n", brightness, blinkSpeed)
}

func (manager DebugHardwareManager) SendBeatToDisplay(force bool) {
	log.Printf("Beat sent, forced: %t \n", force)
}
