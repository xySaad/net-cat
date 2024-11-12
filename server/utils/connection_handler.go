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
	filename := groupName + "_" + time.Now().Format(time.DateOnly) + ".chat.log"
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
	name := string(nameB)
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
	nameb, ok := comands(conn, &name, msg, groupName)
	if ok {
		if nameb != "" {
			name = nameb
		}
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
	filename := groupName + "_" + time.Now().Format(time.DateOnly) + ".chat.log"
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

func comands(conn *net.Conn, name *string, msg []byte, groupName string) (string, bool) {
	if len(msg) != 2 {
		return "", false
	}
	if msg[0] == 8 {
		(*conn).Write([]byte(modules.Comands))
		(*conn).Write(getPrefix((*name)))
		return "", true
	} else if msg[0] == 14 {
		(*conn).Write([]byte("enter your new name: "))
		newName, err := readInput(conn)
		if err != nil {
			(*conn).Write([]byte("an err has occured while changing name"))
		}
		sts := changeName((*name), string(newName), groupName, conn)
		if sts == 0 {
			return string(newName), true
		}
		return "", true
	} else if msg[0] == 12 {
		(*conn).Write([]byte("\033[2J\033[3J\033[H"))
		(*conn).Write(getPrefix((*name)))
		return "", true
	} else if msg[0] == 15 {
		(*conn).Write([]byte("current online members in your group:\n"))
		for v := range modules.Groups.List[groupName] {
			(*conn).Write([]byte(v + "\n"))
		}
		(*conn).Write(getPrefix((*name)))
		return "", true
	} else if msg[0] == 5 {
		filename := groupName + "_" + time.Now().Format(time.DateOnly) + ".chat.log"
		_, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o777)
		if err == nil {
			(*conn).Write([]byte("\033[2J\033[3J\033[H"))
			chat, _ := os.ReadFile(filename)
			(*conn).Write(chat)
			(*conn).Write(getPrefix((*name)))
		}
		return "", true
	}
	return "", false
}
