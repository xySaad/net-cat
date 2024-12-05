package modules

import (
	"net"
	"sync"
)

type User struct {
	net.Conn
	GroupName    string
	UserName     string
	Changingname bool
}
type UsersMap map[string](*User)

type SafeUsers struct {
	mu   sync.Mutex
	List UsersMap
}

func (u *SafeUsers) AddUser(name string, conn *User) uint8 {
	u.mu.Lock()
	defer u.mu.Unlock()

	_, exist := u.List[name]
	if exist {
		return 1
	}

	defer func() {
		conn.UserName = name
	}()

	u.List[name] = conn
	return 0
}

func (u *SafeUsers) DeleteUser(name string) {
	u.mu.Lock()
	defer u.mu.Unlock()

	delete(u.List, name)
}

func (u *SafeUsers) Get(name string) *User {
	u.mu.Lock()
	defer u.mu.Unlock()

	v, ok := u.List[name]

	if !ok {
		return nil
	}

	return v
}
