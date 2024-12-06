package modules

import (
	"fmt"
	"io"
	"os"
	_ "unsafe"

	"net-cat/server/utils"
)

func (conn *User) RestoreHistory() {
	defer conn.Write(utils.GetPrefix(conn.Name))

	err := os.MkdirAll("./history/", 0o755)
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

func GetLogsFileName(groupName string) string {
	return "./history/" + groupName + ".chat.log"
}
