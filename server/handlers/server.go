package handlers

import (
	"fmt"
	"net"
	"os"

	"net-cat/modules"
)

func RunServer(adress string) error {
	ln, err := net.Listen("tcp", adress)
	if err != nil {
		return err
	}

	fmt.Println("server running on:", adress)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		go HandleConnection(&modules.User{Conn: conn})
	}
}
