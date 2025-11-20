package multi

import (
	"fmt"
	"minesweeper-infrastructure/network"
	"minesweeper-infrastructure/protocol"
	"net"
	"os"
	"sync"
)

type Session struct {
	connection *network.Connection
	playerId   int
	mutex      sync.Mutex
}

func NewSession(serverAddress string) (*Session, error) {
	netConn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}
	conn := network.NewConnection(netConn)

	return &Session{
		connection: conn,
	}, nil
}

func (c *Session) JoinGame() error {
	message := protocol.Message{
		Type: protocol.Join,
	}

	return c.connection.Send(message)
}

func (c *Session) StartReceiving() {
	for {
		message, err := c.connection.Receive()
		if err != nil {
			fmt.Println("\n서버 연결이 끊어졌습니다.")
			os.Exit(0)
		}

		c.handleMessage(message)
	}
}

func (c *Session) Close() error {
	return c.connection.Close()
}

func (c *Session) handleMessage(message protocol.Message) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	switch message.Type {
	case protocol.Joined:
		c.handleJoined(message)
	}
}

func (c *Session) handleJoined(message protocol.Message) {
	c.playerId = message.PlayerId
}
