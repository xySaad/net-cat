package modules

const Bitri9 = "         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n[ENTER YOUR NAME]:"

const (
	JoinedStatus = iota
	LeftStatus
	NameChangedStatus
)

type Server struct {
	groups groups
	users  users
}

func NewServer() *Server {
	return &Server{
		groups: groups{
			list: groupsMap{},
		},
		users: users{
			list: usersMap{},
		},
	}
}
