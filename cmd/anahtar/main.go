package main

import (
	"bytes"
	"log"
	"net"
	"strings"
)

const PROTOCOL = "tcp"
const PORT = ":7777"

func main() {
	listener, err := net.Listen(PROTOCOL, PORT)

	if err != nil {
		panic("failed to open listener")
	}

	log.Println("listening on", PORT, "using", PROTOCOL)

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic("we failed while listening")
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Println("we do be handling the connection")

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)

	requestBuffer := bytes.NewBuffer(buf[:n])

	requestString := strings.TrimSpace(requestBuffer.String()) // remove line feed

	if err != nil {
		panic("we failed to read the buffer")
	}

	if requestString == "PING" {
		conn.Write([]byte("PONG"))
	} else {
		conn.Write([]byte("OK"))
	}
}
