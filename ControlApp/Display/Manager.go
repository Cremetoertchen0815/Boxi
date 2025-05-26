package Display

import (
	"encoding/binary"
	"fmt"
	"net"
	"os/exec"
	"sync"
)

type ServerManager struct {
	Connections     map[int]*Server
	connectionMutex *sync.Mutex
	ServerConnected <-chan int
}

func ListenForServers(startLocalDisplayServer bool) (*ServerManager, error) {
	//Open listener
	listener, err := net.Listen("tcp", "localhost:621")
	if err != nil {
		return nil, err
	}

	serverConnected := make(chan int)
	manager := ServerManager{
		make(map[int]*Server),
		&sync.Mutex{},
		serverConnected}

	go listenForClients(&manager, listener, serverConnected)

	if startLocalDisplayServer {
		c := exec.Command("Tools/display_server.py")
		err := c.Start()
		if err != nil {
			return nil, err
		}
	}

	return &manager, nil
}

func listenForClients(manager *ServerManager, listener net.Listener, serverConnected chan<- int) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Listener couldn't be created:", err)
			continue
		}

		go handleClient(manager, conn, serverConnected)
	}
}

func handleClient(manager *ServerManager, conn net.Conn, serverConnected chan<- int) {
	defer conn.Close()

	welcomeBuffer := make([]byte, 3)
	if i, err := conn.Read(welcomeBuffer); i != 3 || err != nil {
		fmt.Println("Welcome message could not be received:", err)
		return
	}

	if welcomeBuffer[0] != 0xE6 || welcomeBuffer[1] != 0x21 {
		fmt.Println("Welcome message had bad header.")
		return
	}

	id := int(welcomeBuffer[2])
	animationReceived := make(chan uint64)
	server := Server{id, conn, animationReceived}

	manager.connectionMutex.Lock()
	manager.Connections[id] = &server
	manager.connectionMutex.Unlock()

	defer func() {
		manager.connectionMutex.Lock()
		delete(manager.Connections, id)
		manager.connectionMutex.Unlock()
	}()

	serverConnected <- id

	for {
		messageBuffer := make([]byte, 6)
		i, err := conn.Read(welcomeBuffer)
		if err != nil {
			return
		}

		if i != 6 || messageBuffer[0] != 0x00 || messageBuffer[1] != 0x01 {
			continue
		}

		//As of now, only one message is supported, that is “AnimationReceived”
		animationId := binary.BigEndian.Uint64(messageBuffer[2:])
		animationReceived <- animationId
	}
}
