package server

import (
	"bytes"
	"log"
	"net"
	"strings"

	"github.com/aurum77/anahtar/internal/handler"
	"github.com/aurum77/anahtar/internal/store"
)

func Start(port string, protocol string) {
	store := store.New()
	handler := handler.New(store)

	listener, err := net.Listen(protocol, port)

	if err != nil {
		panic("Failed to open listener")
	}

	log.Println("Listening on", port, "using", protocol)

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic("Failed while listening")
		}
		go handleConnection(conn, handler)
	}
}

func handleConnection(conn net.Conn, handler *handler.Handler) {
	defer closeConnection(conn)

	log.Println("Handling connection")

	buf := make([]byte, 1024) // do 1024 bytes make sense ?

	n, err := conn.Read(buf)

	if err != nil {
		log.Println("Failed to read data from connection")
		return
	}

	requestBuffer := bytes.NewBuffer(buf[:n])

	command := strings.TrimSpace(requestBuffer.String()) // remove line feed

	handlerResponse := handler.HandleCommand(command)

	conn.Write([]byte(handlerResponse.Message + "\n"))

}

func closeConnection(conn net.Conn) {
	log.Println("Closing connection")
	conn.Close()
}
