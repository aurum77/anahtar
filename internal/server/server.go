package server

import (
	"bufio"
	"io"
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

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("Read failure")
			break
		}

		command := strings.TrimSpace(message) // remove line feed

		handlerResponse := handler.HandleCommand(command)

		_, err = conn.Write([]byte(handlerResponse.Message + "\n"))

		if err != nil {
			log.Println("Write failure")
			break

		}
	}
}

func closeConnection(conn net.Conn) {
	log.Println("Closing connection")
	conn.Close()
}
