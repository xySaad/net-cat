package main

import (
	"fmt"
	"net-cat/handlers"
	"os"
)

var (
	Address = "0.0.0.0"
	Port    = "8989"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "too many arguments")
		return
	}
	if len(os.Args) == 2 {
		Port = os.Args[1]
	}

	adress := Address + ":" + Port
	err := handlers.RunServer(adress)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
