package utils

import (
	"net"
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
	return input, nil
}
