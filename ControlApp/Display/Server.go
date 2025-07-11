package Display

import (
	"encoding/binary"
	"errors"
	"math/rand"
	"net"
	"sync"
	"time"
)

type instructionType byte
type ServerDisplay byte

const (
	DoesAnimationExist instructionType = 0x01
	UploadFrame        instructionType = 0x02
	PlayAnimation      instructionType = 0x03
	ShowText           instructionType = 0x04
	SetBrightness      instructionType = 0x05
)

const (
	Boxi1D1          ServerDisplay = 0b00000001
	Boxi1D2          ServerDisplay = 0b00000010
	Boxi2D1          ServerDisplay = 0b00000100
	Boxi2D2          ServerDisplay = 0b00001000
	allLocalDisplays ServerDisplay = 0b00000011
	AllDisplays      ServerDisplay = 0b00001111
)

type Server struct {
	identifier byte
	connection net.Conn
	callbacks  map[uint32]chan<- bool
	writeLock  *sync.Mutex
}

func (server *Server) sendInstructionWithoutCallback(instructionType instructionType, parameter uint16, payload []byte) {
	data := []byte{'y', 'i', 'f', 'f', byte(instructionType), 0, 0, 0, 0}
	data = binary.BigEndian.AppendUint16(data, parameter)
	data = binary.BigEndian.AppendUint32(data, uint32(len(payload)))
	data = append(data, payload...)

	server.writeLock.Lock()
	defer server.writeLock.Unlock()
	_, _ = server.connection.Write(data)
}

func (server *Server) sendInstructionWithCallback(instructionType instructionType, parameter uint16, payload []byte) (bool, error) {
	callbackId := uint32(rand.Int31())
	callbackCh := make(chan bool, 1)
	server.callbacks[callbackId] = callbackCh

	data := []byte{'y', 'i', 'f', 'f', byte(instructionType)}
	data = binary.BigEndian.AppendUint32(data, callbackId)
	data = binary.BigEndian.AppendUint16(data, parameter)
	data = binary.BigEndian.AppendUint32(data, uint32(len(payload)))
	data = append(data, payload...)

	server.writeLock.Lock()
	defer server.writeLock.Unlock()
	i, err := server.connection.Write(data)

	if i != len(data) {
		return false, errors.New("data wasn't completely transmitted")
	}

	if err != nil {
		return false, err
	}

	result := waitForTimeout(callbackCh)
	delete(server.callbacks, callbackId)
	close(callbackCh)

	return result, nil
}

func waitForTimeout(callback <-chan bool) bool {

	resultCh := make(chan bool, 1)
	go func(callback <-chan bool) {
		select {
		case ret := <-callback:
			resultCh <- ret
		case <-time.After(3 * time.Second):
			resultCh <- false
		}
	}(callback)
	return <-resultCh
}
