package server

import "net"

func HandleConnection(conn *net.Conn) {
	(*conn).Write([]byte(Bitri9))
}
