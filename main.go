package main

import (
	"fmt"
	"net-cat/server"
	"os"
)

var (
	Address = "0.0.0.0"
	Port    = "2000"
)

func main() {
	adress := Address + ":" + Port
	err := server.Run(adress)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
