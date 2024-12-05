package handlers

import (
	"fmt"
	"net"
	"net-cat/server/modules"
	"net-cat/server/utils"
	"os"
)

type Server struct {
	groups modules.SafeGroups
	users  modules.SafeUsers
}

func NewServer() *Server {
	return &Server{
		groups: modules.SafeGroups{
			List: modules.Groups{},
		},
		users: modules.SafeUsers{
			List: modules.UsersMap{},
		},
	}
}

func RunServer(adress string) error {
	server := NewServer()
	ln, err := net.Listen("tcp", adress)
	if err != nil {
		return err
	}

	fmt.Println("server running on:", adress)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		go server.HandleConnection(&modules.User{Conn: conn})
	}
}

func (s *Server) HandleConnection(conn *modules.User) {
	conn.Write([]byte("\033[2J\033[3J\033[H"))
	conn.Write([]byte(modules.Bitri9))

	ok := s.Login(conn, 0)
	if !ok {
		conn.Write([]byte("\n[server]: error login"))
		conn.Close()
		return
	}

	conn.Write([]byte("\033[F\033[2K[ENTER YOUR NAME]:" + conn.UserName + "\n"))
	conn.Write(utils.GetPrefix(conn.UserName))

	s.JoinGroup(conn)
	conn.RestoreHistory()

	s.notify(conn, modules.JoinedStatus)
	for {
		err := s.chat(conn)
		if err != nil {
			break
		}
	}
}
