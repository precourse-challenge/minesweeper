package network

import (
	"encoding/json"
	"minesweeper-infrastructure/protocol"
	"net"
	"sync"
)

type Connection struct {
	Conn    net.Conn
	decoder *json.Decoder
	encoder *json.Encoder
	mutex   sync.Mutex
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		Conn:    conn,
		decoder: json.NewDecoder(conn),
		encoder: json.NewEncoder(conn),
	}
}

func (c *Connection) Receive() (protocol.Message, error) {
	var message protocol.Message
	err := c.decoder.Decode(&message)
	return message, err
}

func (c *Connection) Send(message protocol.Message) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.encoder.Encode(message)
}

func (c *Connection) Close() error {
	return c.Conn.Close()
}
