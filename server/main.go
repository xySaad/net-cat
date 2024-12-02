package main

import (
	"fmt"
	"os"

	"net-cat/handlers"
)

const (
	Address = "0.0.0.0"
	Port    = "8989"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "[USAGE]: ./TCPChat $port")
		return
	}
	adress := Address + ":"

	if len(os.Args) == 2 {
		adress += os.Args[1]
	} else {
		adress += Port
	}

	err := handlers.RunServer(adress)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
