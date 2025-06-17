package Infrastructure

import (
	"ControlApp/BoxiBus"
	"ControlApp/Display"
	"github.com/stianeikeland/go-rpio/v4"
	"log"
	"time"
)

type Manager struct {
	displayServers    *Display.ServerManager
	microController   *BoxiBus.CommunicationHub
	brightness        float64
	blinkSpeed        uint16
	animationProvider AnimationProvider
	beatInput         rpio.Pin
}

const soundInputPin = 16

type AnimationProvider interface {
	getAllAnimations() []Display.AnimationId
}

func Initialize() (*Manager, error) {
	connection, err := BoxiBus.ConnectToArduino(19200)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	displays, err := Display.ListenForServers(true)
	if err != nil {
		message := BoxiBus.CreateDisplayStatusUpdate(BoxiBus.DisplayServerFailed, 1)
		_ = connection.Send(message)

		return &Manager{}, err
	}

	pin := rpio.Pin(soundInputPin)
	pin.Input()

	manager := &Manager{
		displayServers:    displays,
		microController:   connection,
		animationProvider: nil,
		beatInput:         pin,
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
		message := BoxiBus.CreateDisplayStatusUpdate(BoxiBus.HostAwake, serverId)
		err := manager.microController.Send(message)
		if err != nil {
			log.Print(err)
		}

		// Sync animations
		for _, animationId := range manager.animationProvider.getAllAnimations() {
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

func (manager Manager) GetBeatState() bool {
	return manager.beatInput.Read() == rpio.High
}

func (manager Manager) GetConnectedDisplays() []Display.ServerDisplay {
	return manager.displayServers.GetConnectedDisplays()
}

func (manager Manager) UploadAnimation(id Display.AnimationId) {
	frames, err := GetAnimationFrames(uint32(id))
	if err != nil {
		return
	}

	if manager.displayServers.UploadAnimation(id, frames, Display.AllDisplays) != nil {
		log.Printf("Uploading imported animation %d to displays failed. \n", id)
	}
}

func (manager Manager) UpdateStatusCode(statusCode BoxiBus.DisplayStatusCode, serverId byte) {
	message := BoxiBus.CreateDisplayStatusUpdate(statusCode, serverId)
	_ = manager.microController.Send(message)
}
