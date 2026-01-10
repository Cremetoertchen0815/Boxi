package BoxiBus

import (
	"fmt"
	"sync"

	"go.bug.st/serial"
)

type MemoryField byte

type MessageBlock []BusMessage

// Arduino's memory fields/actions
const (
	StatusCode             MemoryField = 0x01
	LightingApply          MemoryField = 0x02
	LightingMode           MemoryField = 0x03
	LightingColorShift     MemoryField = 0x04
	LightingSpeed          MemoryField = 0x05
	LightingGeneralPurpose MemoryField = 0x06
	LightingPaletteSize    MemoryField = 0x07
	LightingPaletteA       MemoryField = 0x08
	LightingPaletteB       MemoryField = 0x09
	EnableInternalLights   MemoryField = 0x10
)

type BusMessage struct {
	field   MemoryField
	payload []byte
}

type CommunicationHub struct {
	lock       *sync.Mutex
	connection serial.Port
}

func ConnectToArduino(baudRate int) (*CommunicationHub, error) {

	mode := &serial.Mode{
		BaudRate: baudRate,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open("/dev/ttyAMA0", mode)
	if err != nil {
		return nil, fmt.Errorf("failed to open UART: %w", err)
	}

	sendMutex := sync.Mutex{}

	return &CommunicationHub{
		&sendMutex,
		port,
	}, nil
}

func (hub *CommunicationHub) sendSingleMessage(message BusMessage) error {
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

	_, err := hub.connection.Write(sendBuffer)
	return err
}

func (hub *CommunicationHub) Send(block MessageBlock) error {
	for _, message := range block {
		if err := hub.sendSingleMessage(message); err != nil {
			return err
		}
	}

	return nil
}
