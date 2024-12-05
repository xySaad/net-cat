package handlers

import (
	"net-cat/server/modules"
	"net-cat/server/utils"
)

const CommandList = "\033[1m\033[32mCTRL + L ===> clear the page\nCTRL + N  ===> option to change your name\nCTRL + H ===> shows all available comands \nCTRL + O ==> shows all online members in the group\nCTRL + E ===> restore chat\n\033[0m"

type commandsMap map[uint8]func()

func (s *Server) executeCommand(conn *modules.User, command uint8) bool {
	var comands = commandsMap{
		'H': func() {
			conn.Write([]byte(CommandList))
			conn.Write([]byte(utils.GetPrefix(conn.UserName)))
		},
		'L': func() {
			conn.Write([]byte("\033[2J\033[3J\033[H"))
			conn.Write([]byte(utils.GetPrefix(conn.UserName)))
		},
		'E': func() {
			conn.Write([]byte("\033[2J\033[3J\033[H"))
			conn.RestoreHistory()
		},
		'O': func() {
			groupMembers := ""
			for member := range s.groups.GetGroup(conn.GroupName) {
				groupMembers += member + "\n"
			}
			conn.Write([]byte("online members:\n" + groupMembers))
			conn.Write([]byte(utils.GetPrefix(conn.UserName)))
		},
		'N': func() {
			s.ChangeName(conn, 0)
		},
	}

	execute, ok := comands[command]
	if ok {
		execute()
		return true
	}
	return false
}
