package modules

import (
	"net"
	"sync"
)

type usersMap map[string](*net.Conn)

type SafeUsers struct {
	sync.Mutex
	List usersMap
}

var Users = SafeUsers{List: make(usersMap)}
