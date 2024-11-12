package main

import (
	"fmt"
	"net"
	"os"

	"net-cat/utils"
)

func Run(adress string) error {
	ln, err := net.Listen("tcp", adress)
	if err != nil {
		return err
	}
	fmt.Println("server running on:", adress)
	os.MkdirAll("./logs/", 0o777)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		go utils.HandleConnection(&conn)
	}
}
