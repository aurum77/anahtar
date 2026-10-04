package server

import (
	"bufio"
	"log"
	"net"

	"github.com/aurum77/anahtar/internal/handler"
	"github.com/aurum77/anahtar/internal/resp"
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
		commandArray := resp.Parse(reader)

		if commandArray == nil {
			break
		} 

		response := handler.HandleCommand(commandArray)

		respResponse := resp.Response(response)

		_, err := conn.Write([]byte(respResponse))

		if err != nil {
			log.Println("Failed to write response")
		}

		log.Printf("ok: %t, msg: %s, response: %s", response.Ok, response.Message, respResponse)
	}
}

func closeConnection(conn net.Conn) {
	log.Println("Closing connection")
	conn.Close()
}
