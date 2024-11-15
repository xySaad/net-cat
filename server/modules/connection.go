package modules

import (
	"fmt"
	"io"
	"net"
	"net-cat/utils"
	"os"
	"strings"
)

type Connection struct {
	net.Conn
	GroupName string
	UserName  string
}

func (conn *Connection) RestoreHistory() {
	err := os.MkdirAll("./logs/", 0755)
	if err != nil {
		fmt.Println(err)
	}

	file, err := os.OpenFile(GetLogsFileName(conn.GroupName), os.O_RDONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			fmt.Println(err)
			conn.Write([]byte("cannot access chat history"))
		}
		return
	}

	defer file.Close()
	chatHistory, err := io.ReadAll(file)

	if err != nil {
		fmt.Println(err)
		conn.Write([]byte("cannot restore chat history"))
		return
	}

	conn.Write(chatHistory)
	conn.Write(utils.GetPrefix(conn.UserName))
}

func (conn *Connection) JoinGroup() {
	conn.Write([]byte("\033[G\033[2K[ENTER GROUP NAME]:"))

	groupNameB, err := utils.ReadInput(&conn.Conn)

	if err != nil {
		if err == io.EOF {
			delete(Users.List, conn.UserName)
		}
		conn.Close()
	}

	groupName := string(groupNameB)
	_, ok := Groups.List[groupName]
	if !ok {
		Groups.List[groupName] = make(map[string]*struct{})
	}
	Groups.List[groupName][conn.UserName] = nil
	conn.GroupName = groupName
	conn.Write([]byte("\033]0;" + groupName + "\a"))
}

func GetLogsFileName(groupName string) string {
	return "./logs/" + groupName + ".chat.log"
}

func (conn *Connection) ChangeName() uint8 {
	newNameB, err := utils.ReadInput(&conn.Conn)

	if err != nil {
		return 1
	}

	newName := string(newNameB)
	if Users.List[newName] != nil {
		conn.Write([]byte("name already taken\n"))
		conn.Write([]byte(utils.GetPrefix(conn.UserName)))
		return 1
	}

	delete(Users.List, conn.UserName)
	delete(Groups.List[conn.GroupName], conn.UserName)
	Users.List[newName] = conn
	Groups.List[conn.GroupName][newName] = nil
	// notify(conn.UserName, conn.GroupName, NameChangedStatus, newName)
	return 0
}

func (conn *Connection) Login(attempts uint8) (string, bool) {
	if attempts > 6 {
		conn.Write([]byte("\033[2K\033[Gtoo many attempts"))
		conn.Close()
		return "", false
	}

	nameB, err := utils.ReadInput(&conn.Conn)
	if err != nil {
		if err == io.EOF {
			return "", false
		}
		fmt.Fprintln(os.Stderr, "error reading from:", conn.RemoteAddr().String())
	}

	if len(nameB) == 0 {
		conn.Write([]byte("empty name is invalid\n[ENTER YOUR NAME]:"))
		return conn.Login(attempts + 1)
	}
	name := string(nameB)
	if attempts > 0 {
		conn.Write([]byte("\033[F\033[2K\033[F\033[2K"))
	}
	if len(name) == 0 {
		conn.Write([]byte("empty name is invalid\n[ENTER YOUR NAME]:"))
		return conn.Login(attempts + 1)
	} else {
		if !utils.ValidUsername(name) {
			conn.Write([]byte("the username " + strings.ReplaceAll(name, string(27), "^[") + " is invalid\n[ENTER YOUR NAME]:"))
			return conn.Login(attempts + 1)
		}
	}

	Users.Lock()
	_, ok := Users.List[name]
	Users.Unlock()

	if ok {
		conn.Write([]byte("the username " + name + " already used\n[ENTER YOUR NAME]:"))
		return conn.Login(attempts + 1)
	}
	Users.Lock()
	conn.UserName = name
	Users.List[name] = conn
	Users.Unlock()

	return name, true
}
