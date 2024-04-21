package main

import "log"

const (
	network = "tcp4"
	addr    = "127.0.0.1:4444"
)

func main() {
	s := NewServer(network, addr)
	go s.Send()
	log.Fatal(s.AcceptConnections())
}
