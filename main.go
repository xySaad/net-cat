package main

import (
	"fmt"
	"log"
	"os"

	"net-cat/server/handlers"
	"net-cat/server/utils"
)

const (
	UDP_ON = iota
	CUSTOM_ADDRESS
	USE_CLIENT
)

const (
	Address  = "0.0.0.0"
	Port     = "8989"
	Protocol = "tcp"
)

func main() {
	status, args := utils.ParseFlags()
	if status == -1 {
		return
	}
	protocol := Protocol
	adress := Address + ":" + Port
	if len(args) > 3 {
		log.Fatalln("err")
	} else if len(args) == 3 {
		adress = args[1] + ":" + args[2]
	} else if len(args) == 2 {
		adress = Address + ":" + args[1]
	}
	if status > 1 {
		protocol = "udp"
	}

	if status == 0 || status == 2 {
		err := handlers.RunServer(adress)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	}

	if status == 1 || status == 3 {
		Client(protocol, adress)
	}
}
