package server

import (
	"fmt"
	"io"
	"net"
	"os"
)

func HandleConnection(conn *net.Conn) {
	(*conn).Write([]byte(Bitri9))

	name, ok := login(conn)
	if !ok {
		return
	}

	brodcast(name, []byte(name+" has joined our chat...\n"))
	chat(name, conn)
}

func login(conn *net.Conn) (string, bool) {
	buffer := make([]byte, 140)
	nameB := []byte{}

	for {
		n, err := (*conn).Read(buffer)

		if err != nil {
			if err == io.EOF {
				return "", false
			}
			fmt.Fprintln(os.Stderr, "error reading from:", (*conn).RemoteAddr().String())
			break
		}

		nameB = append(nameB, buffer[:n]...)
		if buffer[n-1] == '\n' {
			break
		}
	}

	name := string(nameB[:len(nameB)-1])
	Users.Lock()
	_, ok := Users.v[name]
	Users.Unlock()

	if ok {
		(*conn).Write([]byte("the username " + name + " already used\n[ENTER YOUR NAME]:"))
		return login(conn)
	}
	Users.Lock()
	Users.v[name] = conn
	Users.Unlock()

	return name, true
}

func chat(name string, conn *net.Conn) {
	buffer := make([]byte, 140)
	msg := []byte{}

	for {
		n, err := (*conn).Read(buffer)
		if err != nil {
			if err == io.EOF {
				delete(Users.v, string(name))
				brodcast(name, []byte(name+" has left our chat...\n"))
				return
			}
			fmt.Fprintln(os.Stderr, "error reading from:", (*conn).RemoteAddr().String())
			break
		}
		msg = append(msg, buffer[:n]...)
		if buffer[n-1] == '\n' {
			break
		}
	}

	if !(len(msg) == 1 && msg[0] == '\n') {
		brodcast(name, msg)
	}
	chat(name, conn)
}

func brodcast(name string, msg []byte) {
	Users.Lock()
	for user, userConn := range Users.v {
		if user != name {
			(*userConn).Write(msg)
		}
	}
	Users.Unlock()
}
