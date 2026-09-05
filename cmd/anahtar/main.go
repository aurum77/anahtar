package main

import (
	"log"
	"net"
)

const PROTOCOL = "tcp"
const PORT = ":7777"

func main() {
	listener, err := net.Listen(PROTOCOL, PORT)

	if err != nil {
		panic("Failed to open listener")
	}

	log.Println("we have created the listener on", PROTOCOL, PORT)

	for {
		log.Println("we do be listening")
		conn, err := listener.Accept()
		if err != nil {
			panic("We fucked up while listening")
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	log.Println("we do be handling the connection")
}
