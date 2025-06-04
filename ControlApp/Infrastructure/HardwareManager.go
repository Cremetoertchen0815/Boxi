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

type LightingInstruction struct {
	BoxiBus.MessageBlock
	Mode BoxiBus.LightingModeId
}

type AnimationInstruction struct {
	animation  Display.AnimationId
	displays   []Display.ServerDisplay
	blinkSpeed uint16
}

type TextInstruction struct {
	text     string
	displays []Display.ServerDisplay
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

func (manager HardwareManager) SendLightingInstruction(instruction LightingInstruction) {
	err := manager.MicroController.Send(instruction.MessageBlock)
	if err != nil {
		log.Print(err)
	}
}

func (manager HardwareManager) SendAnimationInstruction(instruction AnimationInstruction) {
	totalDisplay := 0
	for _, display := range instruction.displays {
		totalDisplay |= int(display)
	}

	oldVal := manager.blinkSpeed
	manager.blinkSpeed = instruction.blinkSpeed
	if instruction.blinkSpeed != oldVal {
		manager.SendBeatToDisplay(true)
	}

	manager.DisplayServers.PlayAnimation(instruction.animation, Display.ServerDisplay(totalDisplay))
}

func (manager HardwareManager) SendTextInstruction(instruction TextInstruction) {
	totalDisplay := 0
	for display := range instruction.displays {
		totalDisplay |= display
	}
	manager.DisplayServers.DisplayText(instruction.text, Display.ServerDisplay(totalDisplay))
}

func (manager HardwareManager) SendBeatToDisplay(force bool) {
	if !force && manager.blinkSpeed == 0 {
		return
	}

	manager.DisplayServers.SetBrightness(manager.brightness, manager.blinkSpeed)
}
