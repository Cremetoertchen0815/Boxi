package Display

import "net"

type Server struct {
	identifier        int
	connection        net.Conn
	animationReceived <-chan uint64
}
