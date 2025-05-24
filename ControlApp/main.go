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

	colorA := BoxiBus.Color{Blue: 255}
	message := BoxiBus.CreateLightingStrobe(colorA, 3, 0, 0)
	err = connection.Send(message)
	if err != nil {
		log.Fatal(err)
	}
}
