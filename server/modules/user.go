package modules

import (
	"net"
	"sync"
)

type User struct {
	net.Conn
	GroupName    string
	Name         string
	Changingname bool
}
type usersMap map[string](*User)

type users struct {
	mu   sync.Mutex
	list usersMap
}

func (s *Server) StoreUser(name string, conn *User) uint8 {
	s.users.mu.Lock()
	defer s.users.mu.Unlock()

	_, exist := s.users.list[name]
	if exist {
		return 1
	}

	defer func() {
		conn.Name = name
		s.Info("connection", conn.RemoteAddr(), "stored as", conn.Name)
	}()

	s.users.list[name] = conn

	return 0
}

func (s *Server) DeleteUser(name string) {
	s.users.mu.Lock()
	defer s.users.mu.Unlock()

	delete(s.users.list, name)
}

func (s *Server) GetUser(name string) *User {
	s.users.mu.Lock()
	defer s.users.mu.Unlock()

	v, ok := s.users.list[name]

	if !ok {
		return nil
	}

	return v
}
