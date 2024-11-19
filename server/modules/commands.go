package modules

import "net-cat/utils"

const CommandList = "\033[1m\033[32mCTRL + L ===> clear the page\nCTRL + N  ===> option to change your name\nCTRL + H ===> shows all available comands \nCTRL + O ==> shows all online members in the group\nCTRL + E ===> restore chat\n\033[0m"

type commandsMap map[uint8]func(conn *User)

var Comands = commandsMap{
	'H': func(conn *User) {
		conn.Write([]byte(CommandList))
		conn.Write([]byte(utils.GetPrefix(conn.UserName)))
	},
	'L': func(conn *User) {
		conn.Write([]byte("\033[2J\033[3J\033[H"))
		conn.Write([]byte(utils.GetPrefix(conn.UserName)))
	},
	'E': func(conn *User) {
		conn.Write([]byte("\033[2J\033[3J\033[H"))
		conn.RestoreHistory()
	},
	'O': func(conn *User) {
		groupMembers := ""
		for member := range Groups.List[conn.GroupName] {
			groupMembers += member + "\n"
		}
		conn.Write([]byte("online members:\n" + groupMembers))
		conn.Write([]byte(utils.GetPrefix(conn.UserName)))
	},
	18: nil,
}
