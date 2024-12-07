package modules

import (
	"os"
	"time"
)

const Bitri9 = "         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n[ENTER YOUR NAME]:"

const (
	JoinedStatus = iota
	LeftStatus
	NameChangedStatus
)

type Server struct {
	groups groups
	users  users
	logger logger
}

func NewServer() *Server {
	err := os.MkdirAll("./logs/", 0744)
	if err != nil {
		panic(err)
	}
	fd, err := os.OpenFile("./logs/"+time.Now().Format(time.DateOnly)+".log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	return &Server{
		logger: logger{writer: fd},
		groups: groups{
			list: groupsMap{},
		},
		users: users{
			list: usersMap{},
		},
	}
}
