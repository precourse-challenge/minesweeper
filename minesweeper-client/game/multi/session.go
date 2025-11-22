package multi

import (
	"fmt"
	"minesweeper-client/game/view"
	"minesweeper-infrastructure/dto"
	"minesweeper-infrastructure/network"
	"minesweeper-infrastructure/protocol"
	"net"
	"os"
	"sync"
)

type Session struct {
	connection *network.Connection
	playerId   int
	board1     dto.BoardDto
	board2     dto.BoardDto
	gameOver   bool
	mutex      sync.Mutex
	joinDone   chan struct{}
}

func NewSession(serverAddress string) (*Session, error) {
	netConn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}
	conn := network.NewConnection(netConn)

	return &Session{
		connection: conn,
		gameOver:   false,
		joinDone:   make(chan struct{}),
	}, nil
}

func (c *Session) JoinGame() error {
	message := protocol.Message{
		Type: protocol.Join,
	}

	return c.connection.Send(message)
}

func (c *Session) Open(row, col int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.gameOver {
		return nil
	}

	message := protocol.Message{
		Type: protocol.Open,
		Row:  row,
		Col:  col,
	}
	return c.connection.Send(message)
}

func (c *Session) Flag(row, col int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.gameOver {
		return nil
	}

	message := protocol.Message{
		Type: protocol.Flag,
		Row:  row,
		Col:  col,
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
	case protocol.Start:
		c.handleStart(message)
	case protocol.Update:
		c.handleUpdate(message)
	case protocol.Error:
		c.handleError(message)
	}
}

func (c *Session) handleJoined(message protocol.Message) {
	c.playerId = message.PlayerId
	close(c.joinDone)
}

func (c *Session) handleStart(message protocol.Message) {
	c.board1 = message.Board1
	c.board2 = message.Board2
	c.gameOver = false

	view.ShowMultiBoards(c.board1, c.board2, c.playerId)
	view.AskCommand()
}

func (c *Session) handleUpdate(message protocol.Message) {
	c.board1 = message.Board1
	c.board2 = message.Board2

	view.ShowMultiBoards(c.board1, c.board2, c.playerId)
	view.AskCommand()
}

func (c *Session) handleError(message protocol.Message) {
	view.ShowErrorMessage(fmt.Errorf(message.Message))
	view.AskCommand()
}
