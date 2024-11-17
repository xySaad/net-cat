package modules

import (
	"sync"
)

type usersMap map[string](*User)

type SafeUsers struct {
	sync.Mutex
	List usersMap
}

func (u *SafeUsers) AddUser(name string, conn *User) uint8 {
	u.Lock()
	defer u.Unlock()

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
	u.Lock()
	defer u.Unlock()

	delete(u.List, name)
}

func (u *SafeUsers) Get(name string) *User {
	u.Lock()
	defer u.Unlock()

	v, ok := u.List[name]

	if !ok {
		return nil
	}

	return v
}

var Users = SafeUsers{List: make(usersMap)}
