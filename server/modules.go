package server

import (
	"net"
	"sync"
)

var Bitri9 = "         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n[ENTER YOUR NAME]:"

type usersMap map[string](*net.Conn)

type SafeUsers struct {
	sync.Mutex
	list        usersMap
	lastMessage struct {
		msg    string
		sender string
	}
}

var Users = SafeUsers{list: make(usersMap)}
