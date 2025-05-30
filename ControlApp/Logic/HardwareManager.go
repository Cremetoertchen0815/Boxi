package Logic

import (
	"ControlApp/BoxiBus"
	"ControlApp/Display"
	"log"
)

type HardwareManager struct {
	DisplayServers  *Display.ServerManager
	MicroController *BoxiBus.CommunicationHub
}

func InitializeHardware() (HardwareManager, error) {
	//connection, err := BoxiBus.ConnectToArduino(19200)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	defer connection.Close()
	//

	displays, err := Display.ListenForServers(true)
	if err != nil {
		return HardwareManager{}, err
	}

	return HardwareManager{
		DisplayServers:  displays,
		MicroController: nil,
	}, nil
}

func transmitDisplayServerLogon(logonChannel <-chan byte, boxiBus *BoxiBus.CommunicationHub) {
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
	}
}
