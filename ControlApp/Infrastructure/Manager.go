package Infrastructure

import (
	"ControlApp/BoxiBus"
	"ControlApp/Display"
	"log"
	"time"
)

type Manager struct {
	displayServers    *Display.ServerManager
	microController   *BoxiBus.CommunicationHub
	brightness        float64
	blinkSpeed        uint16
	animationProvider AnimationProvider
}

type AnimationProvider interface {
	GetAllAnimations() []Display.AnimationId
}

func Initialize() (Manager, error) {
	connection, err := BoxiBus.ConnectToArduino(19200)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	displays, err := Display.ListenForServers(true)
	if err != nil {
		return Manager{}, err
	}

	manager := Manager{
		displayServers:    displays,
		microController:   connection,
		animationProvider: nil,
	}

	go manager.handleDisplayServerLogon(displays.ServerConnected)

	return manager, nil
}

// handleDisplayServerLogon reports the logon of a display server to the µCs.
func (manager Manager) handleDisplayServerLogon(logonChannel <-chan byte) {
	for {
		if manager.animationProvider == nil {
			time.Sleep(time.Second)
			continue
		}

		serverId := <-logonChannel

		// Send status update to Arduino
		message := BoxiBus.CreateDisplayStatusUpdate(BoxiBus.Active, serverId)
		err := manager.microController.Send(message)
		if err != nil {
			log.Print(err)
		}

		// Sync animations
		for _, animationId := range manager.animationProvider.GetAllAnimations() {
			frames, err := GetAnimationFrames(uint32(animationId))
			if err != nil {
				continue
			}

			_ = manager.displayServers.UploadAnimation(animationId, frames, Display.Boxi1D1<<(serverId*2))
		}
	}
}

func (manager Manager) SendLightingInstruction(block BoxiBus.MessageBlock) {
	err := manager.microController.Send(block)
	if err != nil {
		log.Print(err)
	}
}

func (manager Manager) SendAnimationInstruction(animation Display.AnimationId, displays []Display.ServerDisplay) {
	totalDisplay := 0
	for _, display := range displays {
		totalDisplay |= int(display)
	}

	manager.displayServers.PlayAnimation(animation, Display.ServerDisplay(totalDisplay))
}

func (manager Manager) SendTextInstruction(text string, displays []Display.ServerDisplay) {
	totalDisplay := 0
	for display := range displays {
		totalDisplay |= display
	}
	manager.displayServers.DisplayText(text, Display.ServerDisplay(totalDisplay))
}

func (manager Manager) SendBrightnessChange(brightness *float64, blinkSpeed uint16) {
	if brightness != nil {
		manager.brightness = *brightness
	}

	oldVal := manager.blinkSpeed
	manager.blinkSpeed = blinkSpeed
	if blinkSpeed != oldVal {
		manager.SendBeatToDisplay(true)
	}
}

func (manager Manager) SendBeatToDisplay(force bool) {
	if !force && manager.blinkSpeed == 0 {
		return
	}

	manager.displayServers.SetBrightness(manager.brightness, manager.blinkSpeed)
}
func (manager Manager) GetConnectedDisplays() []Display.ServerDisplay {
	return manager.displayServers.GetConnectedDisplays()
}

func (manager Manager) UploadAnimation(id Display.AnimationId) {
	frames, err := GetAnimationFrames(uint32(id))
	if err != nil {
		return
	}

	manager.displayServers.UploadAnimation(id, frames, Display.AllDisplays)
}
