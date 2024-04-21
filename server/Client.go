package server

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	ID         int
	Connection net.Conn
	Incoming   chan []byte
	Outgoing   chan []byte
	Server     *Server
}

func NewClient(id int, conn net.Conn) *Client {
	ret := new(Client)
	ret.ID = id
	ret.Connection = conn
	ret.Incoming = make(chan []byte)
	ret.Outgoing = make(chan []byte)
	return ret
}

// The Receive method reads in a null-terminated buffer
// from the Client's TCP Connection, and sends the buffer
// to the Incoming channel
func (c *Client) Receive() error {
	defer c.Connection.Close()
	br := bufio.NewReader(c.Connection)
	for {
		buf, err := br.ReadBytes(0)
		if err != nil {
			return err
		}
		br.Reset(c.Connection)
		c.Incoming <- buf
	}
}

// The Send method looks continuously for contents
// the Client's Outgoing channel, and writes them
// to the Client's TCP connection.
func (c *Client) Send() error {
	defer c.Connection.Close()
	for buf := range c.Outgoing {
		_, err := c.Connection.Write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

// The Forward method looks continuously for contents
// in the Client's Outgoing channel, and forwards them
// to the Server's broadcast channel
// func (c *Client) Forward() {
// 	for buf := range c.Incoming {
// 		c.Server.Broadcast <- buf
// 	}
// }

func (c *Client) HandleMessage() {
	for buf := range c.Incoming {
		fmt.Println(buf)
		str := string(buf[:len(buf)-1])
		switch str {
		case "meow":
			fmt.Println("CHEAT CODE ACTIVATED")
		default:
			c.Server.Broadcast <- buf
		}
	}
}
