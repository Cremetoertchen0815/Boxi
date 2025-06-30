package Display

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/exec"
	"sync"
)

type ServerManager struct {
	connections     map[byte]*Server
	connectionMutex *sync.Mutex
	ServerConnected <-chan byte
}

func ListenForServers(startLocalDisplayServer bool) (*ServerManager, error) {
	//Open listener
	listener, err := net.Listen("tcp", "192.168.4.1:25621")
	if err != nil {
		return nil, err
	}

	serverConnected := make(chan byte, 2)
	manager := ServerManager{
		make(map[byte]*Server),
		&sync.Mutex{},
		serverConnected}

	go manager.listenForClients(listener, serverConnected)

	if startLocalDisplayServer {
		c := exec.Command("sudo", "python3", "Display/Server/display_server.py")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err := c.Start()
		if err != nil {
			return nil, err
		}
	}

	return &manager, nil
}

func (manager *ServerManager) GetConnectedDisplays() []ServerDisplay {
	keys := make([]ServerDisplay, 0, len(manager.connections))
	for _, s := range manager.connections {
		keys = append(keys, ServerDisplay(s.identifier*2+1), ServerDisplay(s.identifier*2+2))
	}
	return keys
}

func (manager *ServerManager) listenForClients(listener net.Listener, serverConnected chan<- byte) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Listener couldn't be created:", err)
			continue
		}

		go manager.handleClient(conn, serverConnected)
	}
}

func (manager *ServerManager) handleClient(conn net.Conn, serverConnected chan<- byte) {
	defer conn.Close()

	welcomeBuffer := make([]byte, 7)
	if i, err := conn.Read(welcomeBuffer); i != 7 || err != nil {
		fmt.Println("Welcome message could not be received:", err)
		return
	}

	if welcomeBuffer[0] != 'h' ||
		welcomeBuffer[1] != 'e' ||
		welcomeBuffer[2] != 'w' ||
		welcomeBuffer[3] != 'w' ||
		welcomeBuffer[4] != 'o' ||
		welcomeBuffer[5] != ':' {
		fmt.Println("Welcome message had bad header.")
		return
	}

	id := welcomeBuffer[6]
	callbacks := make(map[uint32]chan<- bool)
	writeLock := sync.Mutex{}
	server := Server{id, conn, callbacks, &writeLock}

	manager.connectionMutex.Lock()
	manager.connections[id] = &server
	manager.connectionMutex.Unlock()

	defer func() {
		manager.connectionMutex.Lock()
		delete(manager.connections, id)
		manager.connectionMutex.Unlock()
	}()

	select {
	case serverConnected <- id:
		// message sent
	default:
		// message dropped
	}

	for {
		messageBuffer := make([]byte, 7)
		i, err := conn.Read(messageBuffer)
		if err != nil {
			return
		}

		if i != 7 || messageBuffer[0] != 0xE6 || messageBuffer[1] != 0x21 {
			continue
		}

		//Report callbacks
		callbackId := binary.BigEndian.Uint32(messageBuffer[2:6])
		callbackCh := callbacks[callbackId]
		callbackCh <- messageBuffer[6] != 0
	}
}
