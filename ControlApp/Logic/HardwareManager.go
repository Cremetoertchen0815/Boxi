package Logic

import (
	"ControlApp/BoxiBus"
	"ControlApp/Display"
	"log"
)

type HardwareManager struct {
	DisplayServers   *Display.ServerManager
	MicroController  *BoxiBus.CommunicationHub
	lightingManual   bool
	animationsManual bool
	textManual       bool
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

		if serverId != 1 {
			continue
		}

		message := BoxiBus.CreateDisplayStatusUpdate(BoxiBus.Active)
		err := boxiBus.Send(message)
		if err != nil {
			log.Print(err)
		}

		// TODO: Sync animations with client when logging on
	}
}
