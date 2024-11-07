package utils

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"net-cat/modules"
)

var (
	joinedStatus = "joined"
	leftStatus   = "left"
)

func HandleConnection(conn *net.Conn) {
	(*conn).Write([]byte("\033[2J\033[3J\033[H"))
	(*conn).Write([]byte(modules.Bitri9))

	name, ok := login(conn, 0)
	if !ok {
		(*conn).Write([]byte("\n[server]: error login"))
		(*conn).Close()
		return
	}
	(*conn).Write([]byte("\033[F\033[2K[ENTER YOUR NAME]:" + name + "\n"))
	(*conn).Write(getPrefix(name))
	greeting(name, joinedStatus)
	chat(name, conn)
}

func login(conn *net.Conn, attempts int) (string, bool) {
	if attempts > 6 {
		(*conn).Write([]byte("\033[2K\033[Gtoo many attempts"))
		(*conn).Close()
		return "", false
	}

	nameB, err := readInput(conn)
	if err != nil {
		if err == io.EOF {
			return "", false
		}
		fmt.Fprintln(os.Stderr, "error reading from:", (*conn).RemoteAddr().String())
	}

	if len(nameB) == 0 {
		(*conn).Write([]byte("empty name is invalid\n[ENTER YOUR NAME]:"))
		return login(conn, attempts+1)
	}
	name := string(nameB[:len(nameB)-1])
	if attempts > 0 {
		(*conn).Write([]byte("\033[F\033[2K\033[F\033[2K"))
	}
	if len(name) == 0 {
		(*conn).Write([]byte("empty name is invalid\n[ENTER YOUR NAME]:"))
		return login(conn, attempts+1)
	} else {
		if !validUsername(name) {
			(*conn).Write([]byte("the username " + strings.ReplaceAll(name, string(27), "^[") + " is invalid\n[ENTER YOUR NAME]:"))
			return login(conn, attempts+1)
		}
	}

	modules.Users.Lock()
	_, ok := modules.Users.List[name]
	modules.Users.Unlock()

	if ok {
		(*conn).Write([]byte("the username " + name + " already used\n[ENTER YOUR NAME]:"))
		return login(conn, attempts+1)
	}
	modules.Users.Lock()
	modules.Users.List[name] = conn
	modules.Users.Unlock()

	return name, true
}

func chat(name string, conn *net.Conn) {
	msg, err := readInput(conn)
	if err != nil {
		if err == io.EOF {
			delete(modules.Users.List, string(name))
			greeting(name, leftStatus)
			return
		}
		fmt.Fprintln(os.Stderr, "error reading from:", (*conn).RemoteAddr().String())
	}

	if !(len(msg) == 1 && msg[0] == '\n') {
		brodcast(name, msg, true)
	} else {
		(*conn).Write([]byte("\033[F\033[2K"))
		(*conn).Write([]byte(getPrefix(name)))
	}

	chat(name, conn)
}

func brodcast(name string, msg []byte, msgPrefix bool) {
	modules.Users.Lock()
	if msgPrefix && !validMsg(msg) {
		(*modules.Users.List[name]).Write([]byte("\033[F\033[2K"))
		(*modules.Users.List[name]).Write([]byte("invalid msg\n"))
		(*modules.Users.List[name]).Write(getPrefix(name))
		modules.Users.Unlock()
		return
	}
	for user, userConn := range modules.Users.List {
		if msgPrefix {
			if user != name {
				(*userConn).Write([]byte("\033[s"))
				(*userConn).Write([]byte{'\n'})
				(*userConn).Write([]byte("\033[F\033[2K"))
			}
			(*userConn).Write(getPrefix(name))
		}
		if user != name {
			if !msgPrefix {
				(*userConn).Write([]byte{'\n'})
				(*userConn).Write([]byte("\033[F\033[2K"))
			}
			(*userConn).Write(msg)
			(*userConn).Write(getPrefix(user))
			if msgPrefix {
				(*userConn).Write([]byte("\033[u\033[B"))
			}
		}
	}
	modules.Users.Unlock()
}

func greeting(name, status string) {
	var msg []byte

	if status == leftStatus {
		msg = []byte(name + " has left our chat...\n")
	} else if status == joinedStatus {
		msg = []byte(name + " has joined our chat...\n")
	}

	brodcast(name, msg, false)
}

func validUsername(name string) bool {
	for _, char := range name {
		if char == 27 {
			return false
		}
		if (!(char >= 'a' && char <= 'z') && !(char >= 'A' && char <= 'z')) && !(char >= '0' && char <= '9') {
			return false
		}
	}
	return true
}

func validMsg(message []byte) bool {
	if len(message) == 0 {
		return false
	}
	for _, char := range string(message) {
		if char == 27 {
			return false
		}
	}
	return true
}

func getPrefix(name string) []byte {
	return []byte("[" + time.Now().Format(time.DateTime) + "][" + name + "]:")
}

func readInput(conn *net.Conn) ([]byte, error) {
	buffer := make([]byte, 140)
	input := []byte{}

	for {
		n, err := (*conn).Read(buffer)
		if err != nil {
			return nil, err
		}
		input = append(input, buffer[:n]...)
		if buffer[n-1] == '\n' {
			break
		}
	}
	return input, nil
}