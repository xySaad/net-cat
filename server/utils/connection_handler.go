package utils

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"net-cat/modules"
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

	groupName, err := joinGroup(name, conn)
	if err != nil {
		(*conn).Close()
		return
	}
	(*conn).Write(getPrefix(name))

	greeting(name, groupName, modules.JoinedStatus)
	chat(name, groupName, conn)
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

func chat(name, groupName string, conn *net.Conn) {
	msg, err := readInput(conn)
	if err != nil {
		if err == io.EOF {
			delete(modules.Users.List, name)
			delete(modules.Groups.List[groupName], name)
			greeting(name, groupName, modules.LeftStatus)
			return
		}
		fmt.Fprintln(os.Stderr, "error reading from:", (*conn).RemoteAddr().String())
	}

	if !(len(msg) == 1 && msg[0] == '\n') {
		brodcast(name, groupName, msg, true)
	} else {
		(*conn).Write([]byte("\033[F\033[2K"))
		(*conn).Write([]byte(getPrefix(name)))
	}

	chat(name, groupName, conn)
}

func brodcast(name, groupName string, msg []byte, msgPrefix bool) {
	modules.Users.Lock()
	if msgPrefix && !validMsg(msg) {
		(*modules.Users.List[name]).Write([]byte("\033[F\033[2K"))
		(*modules.Users.List[name]).Write([]byte("invalid msg\n"))
		(*modules.Users.List[name]).Write(getPrefix(name))
		modules.Users.Unlock()
		return
	}
	modules.Groups.Lock()
	for userName := range modules.Groups.List[groupName] {
		userConn, ok := modules.Users.List[userName]
		if !ok {
			fmt.Println(userName, "is not in the group anymore")
			continue
		}
		if msgPrefix {
			if userName != name {
				(*userConn).Write([]byte("\033[s"))
				(*userConn).Write([]byte{'\n'})
				(*userConn).Write([]byte("\033[F\033[2K"))
			}
			(*userConn).Write(getPrefix(name))
		}
		if userName != name {
			if !msgPrefix {
				(*userConn).Write([]byte{'\n'})
				(*userConn).Write([]byte("\033[F\033[2K"))
			}
			(*userConn).Write(msg)
			(*userConn).Write(getPrefix(userName))
			if msgPrefix {
				(*userConn).Write([]byte("\033[u\033[B"))
			}
		}
	}
	modules.Groups.Unlock()
	modules.Users.Unlock()
}

func greeting(name, groupName, status string) {
	var msg []byte

	if status == modules.LeftStatus {
		msg = []byte(name + " has left our chat...\n")
	} else if status == modules.JoinedStatus {
		msg = []byte(name + " has joined our chat...\n")
	}

	brodcast(name, groupName, msg, false)
}

func joinGroup(name string, conn *net.Conn) (string, error) {
	(*conn).Write([]byte("\033[G\033[2K[ENTER GROUP NAME]:"))
	groupNameB, err := readInput(conn)
	if err != nil {
		return "", err
	}
	groupName := string(groupNameB)
	_, ok := modules.Groups.List[groupName]
	if !ok {
		modules.Groups.List[groupName] = make(map[string]*struct{})
	}
	modules.Groups.List[groupName][name] = nil
	return groupName, nil
}
