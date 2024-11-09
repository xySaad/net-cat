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
	filename := groupName[:len(groupName)-1] + "_" + time.Now().Format(time.DateOnly) + ".chat"
	_, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err == nil {
		chat, _ := os.ReadFile(filename)
		(*conn).Write(chat)
		(*conn).Write(getPrefix(name))
	}
	notify(name, groupName, modules.JoinedStatus)

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
			notify(name, groupName, modules.LeftStatus)
			return
		}
		fmt.Fprintln(os.Stderr, "error reading from:", (*conn).RemoteAddr().String())
	}
	if len(msg) > 2 && strings.HasPrefix(string(msg), string(14)) && validUsername(string(msg[1:len(msg)-1])) {
		sts := changeName(name, string(msg[1:len(msg)-1]), groupName, conn)
		if sts == 0 {
			name = string(msg[1 : len(msg)-1])
		}
	} else if len(msg) == 2 && msg[0] == 12 {
		(*conn).Write([]byte("\033[2J\033[3J\033[H"))
		(*conn).Write(getPrefix(name))
	} else if !(len(msg) == 1 && msg[0] == '\n') {
		brodcast(name, groupName, msg, true)
	} else {
		(*conn).Write([]byte("\033[F\033[2K"))
		(*conn).Write([]byte(getPrefix(name)))
	}
	chat(name, groupName, conn)
}

func changeName(oldName, newName, groupName string, conn *net.Conn) int {
	if modules.Users.List[newName] != nil {
		(*conn).Write([]byte("name already taken\n"))
		(*conn).Write([]byte(getPrefix(oldName)))
		return 1
	}
	delete(modules.Users.List, oldName)
	delete(modules.Groups.List[groupName], oldName)
	modules.Users.List[newName] = conn
	modules.Groups.List[groupName][newName] = nil
	notify(oldName, groupName, modules.NameChangedStatus, newName)
	return 0
}

func brodcast(name, groupName string, msg []byte, msgPrefix bool) {
	valid := validMsg(msg)
	filename := groupName[:len(groupName)-1] + "_" + time.Now().Format(time.DateOnly) + ".chat"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err == nil && valid {
		if msgPrefix {
			file.Write(getPrefix(name))
		}
		file.Write(msg)
	}
	modules.Users.Lock()
	if msgPrefix && !valid {
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

func notify(name, groupName, status string, extra ...string) {
	var msgStr string

	switch status {
	case modules.JoinedStatus:
		msgStr = name + " has joined our chat..."

	case modules.LeftStatus:
		msgStr = name + " has left our chat..."

	case modules.NameChangedStatus:
		msgStr = name + " has changed his name to "
		if len(extra) > 0 {
			msgStr += extra[0]
		}

	default:
	}

	msg := []byte(msgStr + "\n")
	brodcast(name, groupName, msg, false)
}

func joinGroup(name string, conn *net.Conn) (string, error) {
	(*conn).Write([]byte("\033[G\033[2K[ENTER GROUP NAME]:"))
	groupNameB, err := readInput(conn)
	if err != nil {
		if err == io.EOF {
			delete(modules.Users.List, name)
		}
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
