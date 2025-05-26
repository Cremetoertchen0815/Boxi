package main

import (
	"ControlApp/BoxiBus"
	"ControlApp/Display"
	"log"
)

func main() {
	connection, err := BoxiBus.ConnectToArduino(19200)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	displays, err := Display.ListenForServers(true)
	if err != nil {
		log.Fatal(err)
	}

	go transmitDisplayServerLogon(displays.ServerConnected, connection)
}

func transmitDisplayServerLogon(logonChannel <-chan int, boxiBus *BoxiBus.CommunicationHub) {
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
