package Infrastructure

import (
	"ControlApp/BoxiBus"
	"ControlApp/Display"
	"fmt"
	"log"
	"math"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
	"time"
)

type Manager struct {
	displayServers    *Display.ServerManager
	microController   *BoxiBus.CommunicationHub
	brightness        float64
	blinkSpeed        uint16
	animationProvider AnimationProvider
	beatInput         gpio.PinIO
}

type AnimationProvider interface {
	getAllAnimations() []Display.AnimationId
}

func Initialize() (*Manager, error) {
	connection, err := BoxiBus.ConnectToArduino(19200)
	if err != nil {
		log.Fatal(err)
	}

	displays, err := Display.ListenForServers(true)
	if err != nil {
		message := BoxiBus.CreateDisplayStatusUpdate(BoxiBus.DisplayServerFailed, 1)
		_ = connection.Send(message)

		return &Manager{}, err
	}

	pin, err := initBeatPin()
	if err != nil {
		return &Manager{}, fmt.Errorf("GPIOs could not be initialized: %s", err)
	}

	manager := &Manager{
		displayServers:    displays,
		microController:   connection,
		animationProvider: nil,
		beatInput:         pin,
		brightness:        1,
	}

	go manager.handleDisplayServerLogon(displays.ServerConnected)

	return manager, nil
}

func initBeatPin() (gpio.PinIO, error) {
	// Initialize periph.io
	if _, err := host.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize periph: %s", err)
	}

	// Use rpi.P1_11 (GPIO56 on physical pin 19)
	pin := rpi.P1_16

	// Set as input
	if err := pin.In(gpio.PullDown, gpio.NoEdge); err != nil {
		return nil, fmt.Errorf("failed to set pin as input: %s", err)
	}

	return pin, nil
}

// handleDisplayServerLogon reports the logon of a display server to the µCs.
func (manager *Manager) handleDisplayServerLogon(logonChannel <-chan byte) {
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

func (manager *Manager) SendLightingInstruction(block BoxiBus.MessageBlock) {
	err := manager.microController.Send(block)
	if err != nil {
		log.Printf("Error sending lighting instruction: %s", err)
	}
}

func (manager *Manager) SendAnimationInstruction(animation Display.AnimationId, displays []Display.ServerDisplay) {
	totalDisplay := 0
	for _, display := range displays {
		totalDisplay |= int(display)
	}

	manager.displayServers.PlayAnimation(animation, Display.ServerDisplay(totalDisplay))
}

func (manager *Manager) SendTextInstruction(text string, displays []Display.ServerDisplay) {
	totalDisplay := 0
	for _, display := range displays {
		totalDisplay |= int(display)
	}
	manager.displayServers.DisplayText(text, Display.ServerDisplay(totalDisplay))
}

func (manager *Manager) SendBrightnessChange(brightness *float64, blinkSpeed uint16) {
	oldSpeed := manager.blinkSpeed
	oldBrightness := manager.brightness

	if brightness != nil {
		manager.brightness = *brightness
	}

	manager.blinkSpeed = blinkSpeed
	if manager.blinkSpeed != oldSpeed || math.Abs(manager.brightness-oldBrightness) > 0.001 {
		manager.SendBeatToDisplay(true)
	}
}

func (manager *Manager) SendBeatToDisplay(force bool) {
	if !force && manager.blinkSpeed == 0 {
		return
	}
	manager.displayServers.SetBrightness(manager.brightness, manager.blinkSpeed)
}

func (manager *Manager) GetBeatState() bool {
	return manager.beatInput.Read() == gpio.High
}

func (manager *Manager) GetConnectedDisplays() []Display.ServerDisplay {
	return manager.displayServers.GetConnectedDisplays()
}

func (manager *Manager) UploadAnimation(id Display.AnimationId) {
	frames, err := GetAnimationFrames(uint32(id))
	if err != nil {
		return
	}

	if manager.displayServers.UploadAnimation(id, frames, Display.AllDisplays) != nil {
		log.Printf("Uploading imported animation %d to displays failed. \n", id)
	}
}

func (manager *Manager) UpdateStatusCode(statusCode BoxiBus.DisplayStatusCode, serverId byte) {
	message := BoxiBus.CreateDisplayStatusUpdate(statusCode, serverId)
	_ = manager.microController.Send(message)
}
