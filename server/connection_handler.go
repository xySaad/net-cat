package server

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

var (
	joinedStatus = "joined"
	leftStatus   = "left"
)

func HandleConnection(conn *net.Conn) {
	(*conn).Write([]byte("\033[2J\033[3J\033[H"))
	(*conn).Write([]byte(Bitri9))

	name, ok := login(conn, 0)
	if !ok {
		(*conn).Write([]byte("\n[server]: error login"))
		(*conn).Close()
		return
	}

	greeting(name, joinedStatus)
	(*conn).Write([]byte("\033[F[ENTER YOUR NAME]:" + name + "\n"))
	(*conn).Write(getPrefix(name))
	chat(name, conn)
}

func login(conn *net.Conn, attempts int) (string, bool) {
	if attempts > 6 {
		(*conn).Write([]byte("\033[2K\033[Gtoo many attempts"))
		(*conn).Close()
		return "", false
	}
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

	Users.Lock()
	_, ok := Users.list[name]
	Users.Unlock()

	if ok {
		(*conn).Write([]byte("the username " + name + " already used\n[ENTER YOUR NAME]:"))
		return login(conn, attempts+1)
	}
	Users.Lock()
	Users.list[name] = conn
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
				delete(Users.list, string(name))
				greeting(name, leftStatus)
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
		brodcast(name, msg, true)
	} else {
		(*conn).Write([]byte("\033[F\033[2K"))
		(*conn).Write([]byte(getPrefix(name)))
	}

	chat(name, conn)
}

func brodcast(name string, msg []byte, msgPrefix bool) {
	Users.Lock()
	if msgPrefix && !validMsg(msg) {
		(*Users.list[name]).Write([]byte("\033[F\033[2K"))
		(*Users.list[name]).Write([]byte("invalid msg\n"))
		(*Users.list[name]).Write(getPrefix(name))
		Users.Unlock()
		return
	}
	for user, userConn := range Users.list {
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
	Users.Unlock()
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
