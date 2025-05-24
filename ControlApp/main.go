package main

import (
	"ControlApp/BoxiBus"
	"log"
)

func main() {
	connection, err := BoxiBus.ConnectToArduino(19200)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	color := BoxiBus.Color{Green: 255, Blue: 255}
	message := BoxiBus.CreateLightingSetColor(color, color, 0)
	err = connection.Send(message)
	if err != nil {
		log.Fatal(err)
	}
}
