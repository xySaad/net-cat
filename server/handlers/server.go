package handlers

import (
	"fmt"
	"net"
	"net-cat/server/modules"
	"net-cat/server/utils"
)

type TCPServer struct {
	*modules.Server
}

func RunServer(adress string) error {
	server := &TCPServer{modules.NewServer()}

	ln, err := net.Listen("tcp", adress)
	if err != nil {
		server.Error(err)
		return err
	}

	fmt.Println("server running on:", adress)
	server.Info("server running on:", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			server.Error(err)
			continue
		}
		server.Info(conn.RemoteAddr(), "connection accepted")
		go server.HandleConnection(&modules.User{Conn: conn})
	}
}

func (s *TCPServer) HandleConnection(conn *modules.User) {
	conn.Write([]byte("\033[2J\033[3J\033[H"))
	conn.Write([]byte(modules.Bitri9))

	ok := s.Login(conn, 0)
	if !ok {
		s.Warn(conn.RemoteAddr(), "can't login")
		conn.Write([]byte("\n[server]: error login"))
		conn.Close()
		return
	}

	s.Info(conn.RemoteAddr(), "logged-in as", conn.Name)

	conn.Write([]byte("\033[F\033[2K[ENTER YOUR NAME]:" + conn.Name + "\n"))
	conn.Write(utils.GetPrefix(conn.Name))

	joined := s.JoinGroup(conn)
	if !joined {
		s.Warn(conn.RemoteAddr(), "can't join a group")
		return
	}
	conn.RestoreHistory()

	s.notify(conn, modules.JoinedStatus)
	for {
		err := s.chat(conn)
		if err != nil {
			break
		}
	}
}
