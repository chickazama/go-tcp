package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

var (
	nextID = 1
)

const (
	historyFilePath = "./history.txt"
)

type Server struct {
	fp        *os.File
	mtx       sync.Mutex
	Clients   map[int]*Client
	Listener  net.Listener
	Broadcast chan []byte
}

func NewServer(network, addr string) *Server {
	var err error
	ret := new(Server)
	ret.fp, err = os.OpenFile(historyFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	ret.Clients = make(map[int]*Client)
	ret.Broadcast = make(chan []byte, 8)
	ret.Listener, err = net.Listen(network, addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	return ret
}

func (s *Server) Start() {
	go s.AcceptConnections()
	go s.Send()
}

func (s *Server) AcceptConnections() error {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			return err
		}
		s.mtx.Lock()
		client := NewClient(nextID, conn)
		client.Server = s
		s.Clients[client.ID] = client
		go client.Send()
		go client.Receive()
		go client.HandleMessage()
		nextID++
		s.mtx.Unlock()
	}
}

func (s *Server) Send() {
	for buf := range s.Broadcast {
		s.mtx.Lock()
		str := string(buf[:len(buf)-1])
		fmt.Fprintf(s.fp, "%s\n", str)
		for _, c := range s.Clients {
			c.Outgoing <- buf
		}
		s.mtx.Unlock()
	}
}
