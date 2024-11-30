package handlers

import (
	"fmt"
	"io"
	"os"

	"net-cat/modules"
	"net-cat/utils"
)

func HandleConnection(conn *modules.User) {
	conn.Write([]byte("\033[2J\033[3J\033[H"))
	conn.Write([]byte(modules.Bitri9))

	name, ok := conn.Login(0)
	if !ok {
		conn.Write([]byte("\n[server]: error login"))
		conn.Close()
		return
	}

	conn.Write([]byte("\033[F\033[2K[ENTER YOUR NAME]:" + name + "\n"))
	conn.Write(utils.GetPrefix(name))

	conn.JoinGroup()
	conn.RestoreHistory()

	notify(name, conn.GroupName, modules.JoinedStatus)
	for {
		err := chat(&name, conn.GroupName, conn)
		if err != nil {
			break
		}
	}
}

func chat(Name *string, groupName string, conn *modules.User) error {
	name:=(*Name)
	msg, err := utils.ReadInput(&conn.Conn)
	if err != nil {
		if err == io.EOF {
			modules.Users.DeleteUser(name)
			delete(modules.Groups.List[groupName], name)
			notify(name, groupName, modules.LeftStatus)
		} else {
			fmt.Fprintln(os.Stderr, "error reading from:", conn.RemoteAddr().String())
		}
		return err
	}

	if len(msg) == 0 {
		conn.Write([]byte("\033[F\033[2K"))
		conn.Write([]byte(utils.GetPrefix(name)))
		return nil
	}

	if len(msg) == 1 {
		comand, ok := modules.Comands[msg[0]+64]
		if ok {
			comand(conn, &name)
			return nil
		}
	}

	brodcast(name, groupName, msg, true)
	return nil
}

func brodcast(name, groupName string, msg []byte, msgPrefix bool) {
	valid := utils.ValidMsg(msg)
	file, err := os.OpenFile(modules.GetLogsFileName(groupName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)

	if err == nil && valid {
		if msgPrefix {
			file.Write(utils.GetPrefix(name))
			defer file.Write([]byte{'\n'})
		}
		file.Write(msg)
	}

	modules.Users.Lock()
	defer modules.Users.Unlock()

	if msgPrefix && !valid {
		(*modules.Users.List[name]).Write([]byte("\033[F\033[2Kinvalid msg\n"))
		(*modules.Users.List[name]).Write(utils.GetPrefix(name))
		return
	}

	modules.Groups.Lock()
	defer modules.Groups.Unlock()

	for userName := range modules.Groups.List[groupName] {
		userConn := modules.Users.List[userName]

		if msgPrefix {
			if userName != name {
				(*userConn).Write([]byte("\033[s\n\033[F\033[2K"))
			}

			(*userConn).Write(utils.GetPrefix(name))

		}

		if userName != name {

			if !msgPrefix {
				(*userConn).Write([]byte("\n\033[F\033[2K"))
			}

			(*userConn).Write(msg)

			if msgPrefix {
				(*userConn).Write([]byte{'\n'})
				defer (*userConn).Write([]byte("\033[u\033[B"))
			}
			if userConn.Changingname {
				(*userConn).Write([]byte("Enter your new name: "))
			} else {
				(*userConn).Write(utils.GetPrefix(userName))
			}
		}
	}
}

func notify(name, groupName string, status uint8, extra ...string) {
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
