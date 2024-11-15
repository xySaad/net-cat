package utils

import (
	"fmt"
	"io"
	"net"
	"net-cat/modules"
	"os"
	"time"
)

func getPrefix(name string) []byte {
	return []byte("[" + time.Now().Format(time.DateTime) + "][" + name + "]:")
}

func readInput(conn *net.Conn) ([]byte, error) {
	buffer := make([]byte, 140)
	input := []byte{}

	for {
		n, err := (*conn).Read(buffer)
		if err != nil {
			return nil, err
		}
		input = append(input, buffer[:n]...)
		if buffer[n-1] == '\n' {
			break
		}
	}
	return input[:len(input)-1], nil
}

func comands(conn *net.Conn, name *string, msg []byte, groupName string) (string, bool) {
	if len(msg) != 1 {
		return "", false
	}
	if msg[0] == 8 {
		(*conn).Write([]byte(modules.Commands))
		(*conn).Write(getPrefix((*name)))
		return "", true
	} else if msg[0] == 14 {
		(*conn).Write([]byte("enter your new name: "))
		newName, err := readInput(conn)
		if err != nil {
			(*conn).Write([]byte("an err has occured while changing name"))
		}
		sts := changeName((*name), string(newName), groupName, conn)
		if sts == 0 {
			return string(newName), true
		}
		return "", true
	} else if msg[0] == 12 {
		(*conn).Write([]byte("\033[2J\033[3J\033[H"))
		(*conn).Write(getPrefix((*name)))
		return "", true
	} else if msg[0] == 15 {
		(*conn).Write([]byte("current online members in your group:\n"))
		for v := range modules.Groups.List[groupName] {
			(*conn).Write([]byte(v + "\n"))
		}
		(*conn).Write(getPrefix((*name)))
		return "", true
	} else if msg[0] == 5 {
		err := os.MkdirAll("./logs/", 0755)
		if err != nil {
			fmt.Println(err)
		}

		file, err := os.OpenFile(getLogsFileName(groupName), os.O_RDONLY, 0644)
		if err != nil {
			if os.IsExist(err) {
				fmt.Println(err)
				(*conn).Write([]byte("cannot restore chat history"))
			}
		} else {
			defer file.Close()
			chatHistory, err := io.ReadAll(file)
			if err != nil {
				fmt.Println(err)
				(*conn).Write([]byte("cannot restore chat history"))
			} else {
				(*conn).Write([]byte("\033[2J\033[3J\033[H"))
				(*conn).Write(chatHistory)
				(*conn).Write(getPrefix(*name))
			}
		}
		return "", true
	}
	return "", false
}

func getLogsFileName(groupName string) string {
	return "./logs/" + groupName + ".chat.log"
}
