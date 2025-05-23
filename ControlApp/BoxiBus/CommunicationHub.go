package BoxiBus

import (
	"fmt"
	"log"
	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/uart"
	"periph.io/x/conn/v3/uart/uartreg"
	"periph.io/x/host/v3"
	"sync"
)

type MemoryField byte

type MessageBlock []BusMessage

const (
	StatusCode             MemoryField = 0x01
	LightingApply          MemoryField = 0x02
	LightingMode           MemoryField = 0x03
	LightingColorShift     MemoryField = 0x04
	LightingSpeed          MemoryField = 0x05
	LightingGeneralPurpose MemoryField = 0x06
	LightingPaletteSize    MemoryField = 0x07
	LightingPaletteA       MemoryField = 0x08
)

type BusMessage struct {
	field   MemoryField
	payload []byte
}

type CommunicationHub struct {
	lock       *sync.Mutex
	connection conn.Conn
	portClose  *uart.PortCloser
}

func ConnectToArduino(baudRate int) CommunicationHub {
	// Initialize host
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
		connection,
		&port,
	}
}

func (hub *CommunicationHub) Send(message BusMessage) error {
	payloadLen := len(message.payload)
	if payloadLen > 6 {
		return fmt.Errorf("the payload length cannot exceed 6 bytes, but payload is %d bytes", payloadLen)
	}

	sendBuffer := make([]byte, payloadLen+4)
	sendBuffer[0] = 0x55
	sendBuffer[1] = 0x77
	sendBuffer[2] = 0x4f
	sendBuffer[3] = byte(message.field)
	for i := 0; i < len(message.payload); i++ {
		sendBuffer[4+i] = message.payload[i]
	}

	receiveBuffer := make([]byte, payloadLen+4)
	return hub.connection.Tx(sendBuffer, receiveBuffer)
}

func (hub *CommunicationHub) SendBlock(block MessageBlock) error {
	for _, message := range block {
		if err := hub.Send(message); err != nil {
			return err
		}
	}

	return nil
}
