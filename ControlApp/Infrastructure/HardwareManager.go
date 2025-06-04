package Infrastructure

import (
	"ControlApp/BoxiBus"
	"ControlApp/Display"
	"log"
)

type HardwareManager struct {
	DisplayServers  *Display.ServerManager
	MicroController *BoxiBus.CommunicationHub
	brightness      float64
	blinkSpeed      uint16
}

func InitializeHardware() (HardwareManager, error) {
	connection, err := BoxiBus.ConnectToArduino(19200)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	displays, err := Display.ListenForServers(true)
	if err != nil {
		return HardwareManager{}, err
	}

	go handleDisplayServerLogon(displays.ServerConnected, connection)

	return HardwareManager{
		DisplayServers:  displays,
		MicroController: nil,
	}, nil
}

// handleDisplayServerLogon reports the logon of a display server to the µCs.
func handleDisplayServerLogon(logonChannel <-chan byte, boxiBus *BoxiBus.CommunicationHub) {
	for {
		serverId := <-logonChannel

		message := BoxiBus.CreateDisplayStatusUpdate(BoxiBus.Active, serverId)
		err := boxiBus.Send(message)
		if err != nil {
			log.Print(err)
		}

		// TODO: Sync animations with client when logging on
	}
}

func (manager HardwareManager) SendLightingInstruction(block BoxiBus.MessageBlock) {
	err := manager.MicroController.Send(block)
	if err != nil {
		log.Print(err)
	}
}

func (manager HardwareManager) SendAnimationInstruction(animation Display.AnimationId, displays []Display.ServerDisplay) {
	totalDisplay := 0
	for _, display := range displays {
		totalDisplay |= int(display)
	}

	manager.DisplayServers.PlayAnimation(animation, Display.ServerDisplay(totalDisplay))
}

func (manager HardwareManager) SendTextInstruction(text string, displays []Display.ServerDisplay) {
	totalDisplay := 0
	for display := range displays {
		totalDisplay |= display
	}
	manager.DisplayServers.DisplayText(text, Display.ServerDisplay(totalDisplay))
}

func (manager HardwareManager) SendBrightnessChange(brightness *float64, blinkSpeed uint16) {
	if brightness != nil {
		manager.brightness = *brightness
	}

	oldVal := manager.blinkSpeed
	manager.blinkSpeed = blinkSpeed
	if blinkSpeed != oldVal {
		manager.SendBeatToDisplay(true)
	}
}

func (manager HardwareManager) SendBeatToDisplay(force bool) {
	if !force && manager.blinkSpeed == 0 {
		return
	}

	manager.DisplayServers.SetBrightness(manager.brightness, manager.blinkSpeed)
}
