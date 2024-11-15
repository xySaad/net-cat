package modules

import (
	"sync"
)

type usersMap map[string](*Connection)

type SafeUsers struct {
	sync.Mutex
	List usersMap
}

var Users = SafeUsers{List: make(usersMap)}
