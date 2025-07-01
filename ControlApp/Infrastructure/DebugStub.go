package Infrastructure

import (
	"ControlApp/BoxiBus"
	"ControlApp/Display"
	"log"
)

type DebugStub struct {
	BeatTriggered bool
}

func (manager DebugStub) GetConnectedDisplays() []Display.ServerDisplay {
	return []Display.ServerDisplay{Display.ServerDisplay(1)}
}
func (manager DebugStub) GetBeatState() bool {
	result := manager.BeatTriggered
	manager.BeatTriggered = false
	return result
}

func (manager DebugStub) SendLightingInstruction(block BoxiBus.MessageBlock) {
	log.Printf("Lighting instruction sent: %+v \n", block)
}

func (manager DebugStub) SendAnimationInstruction(animation Display.AnimationId, displays []Display.ServerDisplay) {
	log.Printf("Animation instruction sent, animation: %d, displays: %+v \n", animation, displays)
}

func (manager DebugStub) SendTextInstruction(text string, displays []Display.ServerDisplay) {
	log.Printf("Text instruction sent, text: %s, displays: %+v \n", text, displays)
}

func (manager DebugStub) SendBrightnessChange(brightness *float64, blinkSpeed uint16) {
	log.Printf("Brightness instruction sent, brightness: %+v, speed: %d \n", brightness, blinkSpeed)
}

func (manager DebugStub) SendBeatToDisplay(force bool) {
	log.Printf("Beat sent, forced: %t \n", force)
}

func (manager DebugStub) UploadAnimation(id Display.AnimationId) {
	log.Printf("Uploaded animation, id: %d \n", id)
}

func (manager DebugStub) UpdateStatusCode(statusCode BoxiBus.DisplayStatusCode, serverId byte) {
	log.Printf("Updated status code, code: %d, serverId: %d \n", statusCode, serverId)
}
