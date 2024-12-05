package handlers

import (
	"fmt"
	"io"
	"os"
	"strings"

	"net-cat/server/modules"
	"net-cat/server/utils"
)

func (s *Server) chat(conn *modules.User) error {
	msg, err := utils.ReadInput(&conn.Conn)
	if err != nil {
		if err == io.EOF {
			s.users.DeleteUser(conn.UserName)
			delete(s.groups.GetGroup(conn.GroupName), conn.UserName)
			s.notify(conn, modules.LeftStatus)
		} else {
			fmt.Fprintln(os.Stderr, "error reading from:", conn.RemoteAddr().String())
		}
		return err
	}

	if len(msg) == 0 {
		conn.Write([]byte("\033[F\033[2K"))
		conn.Write([]byte(utils.GetPrefix(conn.UserName)))
		return nil
	}

	if len(msg) == 1 {
		valid := s.executeCommand(conn, msg[0]+64)
		if valid {
			return nil
		}
	}

	s.brodcast(conn, msg, true)
	return nil
}

func (s *Server) brodcast(conn *modules.User, msg []byte, msgPrefix bool) {
	valid := utils.ValidMsg(msg)
	file, err := os.OpenFile(modules.GetLogsFileName(conn.GroupName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)

	if err == nil && valid {
		if msgPrefix {
			file.Write(utils.GetPrefix(conn.UserName))
			defer file.Write([]byte{'\n'})
		}
		file.Write(msg)
	}

	if msgPrefix && !valid {
		s.users.Get(conn.UserName).Write([]byte("\033[F\033[2Kinvalid msg\n"))
		s.users.Get(conn.UserName).Write(utils.GetPrefix(conn.UserName))
		return
	}

	for userName := range s.groups.GetGroup(conn.GroupName) {
		userConn := s.users.Get(userName)

		if msgPrefix {
			if userName != conn.UserName {
				userConn.Write([]byte("\033[s\n\033[F\033[2K"))
			}

			userConn.Write(utils.GetPrefix(conn.UserName))

		}

		if userName != conn.UserName {

			if !msgPrefix {
				userConn.Write([]byte("\n\033[F\033[2K"))
			}

			userConn.Write(msg)

			if msgPrefix {
				userConn.Write([]byte{'\n'})
				defer userConn.Write([]byte("\033[u\033[B"))
			}
			if userConn.Changingname {
				userConn.Write([]byte("Enter your new name: "))
			} else {
				userConn.Write(utils.GetPrefix(userName))
			}
		}
	}
}

func (s *Server) notify(conn *modules.User, status uint8, extra ...string) {
	var msgStr string

	switch status {
	case modules.JoinedStatus:
		msgStr = "\033[38;2;0;184;30m" + conn.UserName + " has joined our chat..."

	case modules.LeftStatus:
		msgStr = "\033[38;2;255;0;0m" + conn.UserName + " has left our chat..."

	case modules.NameChangedStatus:
		color := "\033[38;2;146;142;210m"
		msg := " has changed his name to " + conn.UserName
		if len(extra) > 0 {
			msgStr = color + extra[0] + msg
		} else {
			msgStr = color + "someone" + msg
		}
	default:
	}

	msg := []byte(msgStr + "\n\033[0m")
	s.brodcast(conn, msg, false)
}

func (s *Server) JoinGroup(conn *modules.User) {
	conn.Write([]byte("\033[G\033[2K[ENTER GROUP NAME]:"))

	groupNameB, err := utils.ReadInput(&conn.Conn)
	if err != nil {
		if err == io.EOF {
			s.users.DeleteUser(conn.UserName)
		}
		conn.Close()
		return
	}

	groupName := string(groupNameB)
	status := utils.ValidName(groupName)
	if status != 0 {
		if status == 1 {
			conn.Write([]byte("the group name can be at least 3 characters"))
		}
		if status == 2 {
			conn.Write([]byte("the group name cannot be more than 12 characters"))
		}
		if status == 3 {
			conn.Write([]byte("the group name can only contain alphanumerical characters (a-z_0-9)"))
		}
		conn.Write([]byte("\n[ENTER YOUR NAME]:"))
		s.JoinGroup(conn)
		return
	}
	groupName += "_" + strings.Split(conn.Conn.LocalAddr().String(), ":")[1]
	s.groups.AddUser(groupName, conn)
	conn.Write([]byte("\033]0;" + groupName + "\a"))
}

func (s *Server) Login(conn *modules.User, attempts uint8) bool {
	if attempts > 6 {
		conn.Write([]byte("\033[2K\033[Gtoo many attempts"))
		conn.Close()
		return false
	}

	nameB, err := utils.ReadInput(&conn.Conn)
	if err != nil {
		if err != io.EOF {
			fmt.Fprintln(os.Stderr, "error reading from1:", conn.RemoteAddr())
		}
		return false
	}

	name := string(nameB)
	if attempts > 0 {
		conn.Write([]byte("\033[F\033[2K\033[F\033[2K"))
	}
	if len(name) == 0 {
		conn.Write([]byte("empty name is invalid\n[ENTER YOUR NAME]:"))
		return s.Login(conn, attempts+1)
	}

	status := utils.ValidName(name)
	if status != 0 {
		if status == 1 {
			conn.Write([]byte("the username can be at least 3 characters"))
		}
		if status == 2 {
			conn.Write([]byte("the username cannot be more than 12 characters"))
		}
		if status == 3 {
			conn.Write([]byte("the username can only contain alphanumerical characters (a-z_0-9)"))
		}
		conn.Write([]byte("\n[ENTER YOUR NAME]:"))
		return s.Login(conn, attempts+1)
	}

	status = s.users.AddUser(name, conn)
	if status == 1 {
		conn.Write([]byte("the username " + name + " already used\n[ENTER YOUR NAME]:"))
		return s.Login(conn, attempts+1)
	}
	conn.UserName = name

	return true
}

func (s *Server) ChangeName(conn *modules.User, try int) uint8 {
	if try == 5 {
		conn.Write([]byte("too many attempts...\n"))
		conn.Write([]byte(utils.GetPrefix(conn.UserName)))
		conn.Changingname = false
		return 1
	}
	conn.Write([]byte("Enter your new name: "))
	conn.Changingname = true
	newNameB, err := utils.ReadInput(&conn.Conn)
	if err != nil {
		conn.Changingname = false
		return 1
	}

	newName := string(newNameB)
	status := utils.ValidName(newName)
	if status != 0 {
		conn.Write([]byte("\033[F\033[2K\033[2K"))
		if status == 1 {
			conn.Write([]byte("the username can be at least 3 characters\n"))
		}
		if status == 2 {
			conn.Write([]byte("the username cannot be more than 12 characters\n"))
		}
		if status == 3 {
			conn.Write([]byte("the username can only contain alphanumerical characters (a-z_0-9)\n"))
		}
		s.ChangeName(conn, try+1)
		return 0
	}
	if s.users.Get(newName) != nil {
		conn.Write([]byte("name already taken\n"))
		s.ChangeName(conn, try+1)
		return 0
	}

	defer s.notify(conn, modules.NameChangedStatus, conn.UserName)

	s.users.DeleteUser(conn.UserName)
	s.groups.DeleteFromGroup(conn)
	conn.UserName = newName
	s.users.AddUser(newName, conn)
	s.groups.AddUser(conn.GroupName, conn)
	conn.Write([]byte(utils.GetPrefix(conn.UserName)))
	conn.Changingname = false
	return 0
}
