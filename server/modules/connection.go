package modules

import (
	"fmt"
	"io"
	"net"
	"os"
	_ "unsafe"

	"net-cat/utils"
)

type User struct {
	net.Conn
	GroupName    string
	UserName     string
	Changingname bool
}

func (conn *User) RestoreHistory() {
	defer conn.Write(utils.GetPrefix(conn.UserName))

	err := os.MkdirAll("./logs/", 0o755)
	if err != nil {
		fmt.Println(err)
		conn.Write([]byte("cannot access chat history"))
		return
	}

	file, err := os.OpenFile(GetLogsFileName(conn.GroupName), os.O_RDONLY, 0o644)
	if err != nil {
		if os.IsExist(err) {
			fmt.Fprintln(os.Stderr, err)
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
}

func (conn *User) JoinGroup() {
	conn.Write([]byte("\033[G\033[2K[ENTER GROUP NAME]:"))

	groupNameB, err := utils.ReadInput(&conn.Conn)
	if err != nil {
		if err == io.EOF {
			Users.DeleteUser(conn.UserName)
		}
		conn.Close()
		return
	}

	groupName := string(groupNameB)
	status := utils.ValidName(groupName)
	if status != 0 {
		if status == 1 {
			conn.Write([]byte("the group name can be at least 3 characters"))
		}
		if status == 2 {
			conn.Write([]byte("the group name cannot be more than 12 characters"))
		}
		if status == 3 {
			conn.Write([]byte("the group name can only contain alphanumerical characters (a-z_0-9)"))
		}
		conn.Write([]byte("\n[ENTER YOUR NAME]:"))
		conn.JoinGroup()
		return
	}

	Groups.SetGroup(groupName, conn)
	conn.Write([]byte("\033]0;" + groupName + "\a"))
}

func GetLogsFileName(groupName string) string {
	return "./logs/" + groupName + ".chat.log"
}

//go:linkname handlers_notify net-cat/handlers.notify
func handlers_notify(name, groupName string, status uint8, extra ...string)

func (conn *User) ChangeName(name *string, try int) uint8 {
	if try == 5 {
		conn.Write([]byte("too many attempts...\n"))
		conn.Write([]byte(utils.GetPrefix(conn.UserName)))
		conn.Changingname = false
		return 1
	}
	conn.Write([]byte("Enter your new name: "))
	conn.Changingname = true
	newNameB, err := utils.ReadInput(&conn.Conn)
	if err != nil {
		conn.Changingname = false
		return 1
	}

	newName := string(newNameB)
	status := utils.ValidName(newName)
	if status != 0 {
		conn.Write([]byte("\033[F\033[2K\033[2K"))

		if status == 1 {
			conn.Write([]byte("the username can be at least 3 characters\n"))
		}
		if status == 2 {
			conn.Write([]byte("the username cannot be more than 12 characters\n"))
		}
		if status == 3 {
			conn.Write([]byte("the username can only contain alphanumerical characters (a-z_0-9)\n"))
		}
		conn.ChangeName(name, try+1)
		return 0
	}
	if Users.Get(newName) != nil {
		conn.Write([]byte("name already taken\n"))
		conn.ChangeName(name, try+1)
		return 0
	}

	Users.DeleteUser(conn.UserName)
	delete(Groups.List[conn.GroupName], conn.UserName)
	(*name) = newName
	Users.AddUser(newName, conn)
	Groups.List[conn.GroupName][newName] = nil
	handlers_notify(conn.UserName, conn.GroupName, NameChangedStatus, newName)
	conn.Write([]byte(utils.GetPrefix(conn.UserName)))
	return 0
}

func (conn *User) Login(attempts uint8) (string, bool) {
	if attempts > 6 {
		conn.Write([]byte("\033[2K\033[Gtoo many attempts"))
		conn.Close()
		return "", false
	}

	nameB, err := utils.ReadInput(&conn.Conn)
	if err != nil {
		if err != io.EOF {
			fmt.Fprintln(os.Stderr, "error reading from1:", conn.RemoteAddr())
		}
		return "", false
	}

	name := string(nameB)
	if attempts > 0 {
		conn.Write([]byte("\033[F\033[2K\033[F\033[2K"))
	}
	if len(name) == 0 {
		conn.Write([]byte("empty name is invalid\n[ENTER YOUR NAME]:"))
		return conn.Login(attempts + 1)
	}

	status := utils.ValidName(name)
	if status != 0 {
		if status == 1 {
			conn.Write([]byte("the username can be at least 3 characters"))
		}
		if status == 2 {
			conn.Write([]byte("the username cannot be more than 12 characters"))
		}
		if status == 3 {
			conn.Write([]byte("the username can only contain alphanumerical characters (a-z_0-9)"))
		}
		conn.Write([]byte("\n[ENTER YOUR NAME]:"))
		return conn.Login(attempts + 1)
	}

	status = Users.AddUser(name, conn)
	if status == 1 {
		conn.Write([]byte("the username " + name + " already used\n[ENTER YOUR NAME]:"))
		return conn.Login(attempts + 1)
	}

	return name, true
}
