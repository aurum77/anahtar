package server

import (
	"bufio"
	"log"
	"net"
	"strconv"

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
		parseResult := resp.Parse(reader)

		if parseResult.Commands == nil {
			break
		}

		response := handler.HandleCommand(parseResult.Commands)

		respResponse := resp.Response(response)

		byteCount, err := conn.Write([]byte(respResponse))

		if byteCount != len(respResponse) {
			log.Println("Response bytecount doesn't match written bytecount")
		}

		if err != nil {
			log.Println("Failed to write response, closing connection")
			conn.Close()
		}

		log.Printf("ok: %s, msg: %s, response: %s", strconv.Itoa(response.Type), response.Message, respResponse)
	}
}

func closeConnection(conn net.Conn) {
	log.Println("Closing connection")
	conn.Close()
}
