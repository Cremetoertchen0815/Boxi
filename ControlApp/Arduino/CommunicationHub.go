package Arduino

import (
	"log"
	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/uart"
	"periph.io/x/conn/v3/uart/uartreg"
	"periph.io/x/host/v3"
	"sync"
)

type CommunicationHub struct {
	lock       *sync.Mutex
	connection *conn.Conn
	portClose  *uart.PortCloser
}

func ConnectToArduino(baudRate int) CommunicationHub {
	// Initialize periph
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Open UART port
	port, err := uartreg.Open("/dev/serial0")
	if err != nil {
		log.Fatal(err)
	}

	connection, err := port.Connect(physic.Frequency(baudRate), uart.One, uart.NoParity, uart.NoFlow, 8)
	if err != nil {
		log.Fatal(err)
	}

	sendMutex := sync.Mutex{}

	return CommunicationHub{
		&sendMutex,
		&connection,
		&port,
	}
}
