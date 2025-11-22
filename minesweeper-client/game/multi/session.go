package multi

import (
	"fmt"
	"minesweeper-infrastructure/dto"
	"minesweeper-infrastructure/network"
	"minesweeper-infrastructure/protocol"
	"net"
	"os"
	"sync"
)

type Session struct {
	connection    *network.Connection
	playerId      int
	board1        dto.BoardDto
	board2        dto.BoardDto
	gameOver      bool
	mutex         sync.Mutex
	eventChannels *SessionEventChannels
}

func NewSession(serverAddress string, eventChannels *SessionEventChannels) (*Session, error) {
	netConn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}

	return &Session{
		connection:    network.NewConnection(netConn),
		eventChannels: eventChannels,
	}, nil
}

func (c *Session) JoinGame() error {
	return c.connection.Send(protocol.Message{Type: protocol.Join})
}

func (c *Session) Open(row, col int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.gameOver {
		return nil
	}
	return c.connection.Send(protocol.Message{Type: protocol.Open, Row: row, Col: col})
}

func (c *Session) Flag(row, col int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.gameOver {
		return nil
	}
	return c.connection.Send(protocol.Message{Type: protocol.Flag, Row: row, Col: col})
}

func (c *Session) StartReceiving() {
	for {
		msg, err := c.connection.Receive()
		if err != nil {
			fmt.Println("\n서버 연결이 끊어졌습니다.")
			os.Exit(0)
		}
		c.handleMessage(msg)
	}
}

func (c *Session) handleMessage(message protocol.Message) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	switch message.Type {
	case protocol.Joined:
		c.handleJoined(message)
	case protocol.Start:
		c.handleStart(message)
	case protocol.Update:
		c.handleUpdate(message)
	case protocol.Error:
		c.handleError(message)
	case protocol.GameOver:
		c.handleGameOver(message)
	}
}

func (c *Session) handleJoined(message protocol.Message) {
	c.playerId = message.PlayerId
	c.eventChannels.JoinedChan <- JoinedEvent{PlayerId: c.playerId}
}

func (c *Session) handleStart(message protocol.Message) {
	c.board1 = message.Board1
	c.board2 = message.Board2
	c.gameOver = false

	c.eventChannels.StartChan <- StartEvent{
		Board1:   c.board1,
		Board2:   c.board2,
		PlayerId: c.playerId,
	}
}

func (c *Session) handleUpdate(message protocol.Message) {
	c.board1 = message.Board1
	c.board2 = message.Board2

	c.eventChannels.UpdateChan <- UpdateEvent{
		Board1:   c.board1,
		Board2:   c.board2,
		PlayerId: c.playerId,
	}
}

func (c *Session) handleError(message protocol.Message) {
	c.eventChannels.ErrorChan <- ErrorEvent{
		Err: fmt.Errorf(message.Message),
	}
}

func (c *Session) handleGameOver(message protocol.Message) {
	c.board1 = message.Board1
	c.board2 = message.Board2
	c.gameOver = true

	c.eventChannels.GameOverChan <- GameOverEvent{
		Board1:   c.board1,
		Board2:   c.board2,
		PlayerId: c.playerId,
		Winner:   message.Winner,
	}
}

func (c *Session) Close() error {
	return c.connection.Close()
}
