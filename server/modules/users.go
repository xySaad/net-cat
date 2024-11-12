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

var (
	Comands = "CTRL + L ===> clear the page\nCTRL + N  ===> option to change your name\nCTRL + H ===> shows all available comands \nCTRL + O ==> shows all online members in the group\nCTRL + E ===> restore chat\n"
	Users   = SafeUsers{List: make(usersMap)}
)
