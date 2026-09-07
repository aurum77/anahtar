package main

import "github.com/aurum77/anahtar/internal/server"

const PROTOCOL = "tcp"
const PORT = ":7777"

func main() {
	server.Start(PORT, PROTOCOL)
}
