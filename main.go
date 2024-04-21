package main

import (
	"log"

	"github.com/chickazama/go-tcp/server"
)

const (
	network = "tcp4"
	addr    = "127.0.0.1:4444"
)

func main() {
	s := server.NewServer(network, addr)
	go s.Send()
	log.Fatal(s.AcceptConnections())
}
