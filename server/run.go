package server

import (
	"fmt"
	"net"
	"os"
)

func Run(adress string) error {
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
		go HandleConnection(&conn)
	}
}
