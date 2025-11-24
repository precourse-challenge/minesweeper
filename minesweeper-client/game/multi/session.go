package multi

import (
	"fmt"
	"minesweeper-infrastructure/dto"
	"minesweeper-infrastructure/network"
	"minesweeper-infrastructure/protocol"
	"net"
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

func (s *Session) JoinGame() error {
	return s.connection.Send(protocol.Message{Type: protocol.Join})
}

func (s *Session) Open(row, col int) error {
	if s.gameOver {
		return nil
	}
	return s.connection.Send(protocol.Message{
		Type: protocol.Open,
		Row:  row,
		Col:  col,
	})
}

func (s *Session) Flag(row, col int) error {
	if s.gameOver {
		return nil
	}
	return s.connection.Send(protocol.Message{
		Type: protocol.Flag,
		Row:  row,
		Col:  col,
	})
}

func (s *Session) Exit() error {
	return s.connection.Send(protocol.Message{
		Type:     protocol.Exit,
		PlayerId: s.playerId,
	})
}

func (s *Session) StartReceiving() {
	for {
		message, err := s.connection.Receive()
		if err != nil {
			return
		}
		s.handleMessage(message)
	}
}

func (s *Session) handleMessage(message protocol.Message) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	switch message.Type {
	case protocol.Joined:
		s.handleJoined(message)
	case protocol.Start:
		s.handleStart(message)
	case protocol.Update:
		s.handleUpdate(message)
	case protocol.Error:
		s.handleError(message)
	case protocol.GameOver:
		s.handleGameOver(message)
	}
}

func (s *Session) handleJoined(message protocol.Message) {
	s.playerId = message.PlayerId
	s.eventChannels.JoinedChan <- JoinedEvent{
		PlayerId: s.playerId,
		Message:  message.Message,
	}
}

func (s *Session) handleStart(message protocol.Message) {
	s.board1 = message.Board1
	s.board2 = message.Board2
	s.gameOver = false

	s.eventChannels.StartChan <- StartEvent{
		Board1:   s.board1,
		Board2:   s.board2,
		PlayerId: s.playerId,
		Message:  message.Message,
	}
}

func (s *Session) handleUpdate(message protocol.Message) {
	s.board1 = message.Board1
	s.board2 = message.Board2

	s.eventChannels.UpdateChan <- UpdateEvent{
		Board1:   s.board1,
		Board2:   s.board2,
		PlayerId: s.playerId,
		Message:  message.Message,
	}
}

func (s *Session) handleError(message protocol.Message) {
	s.eventChannels.ErrorChan <- ErrorEvent{
		Err:     fmt.Errorf(message.Message),
		Message: message.Message,
	}
}

func (s *Session) handleGameOver(message protocol.Message) {
	s.board1 = message.Board1
	s.board2 = message.Board2
	s.gameOver = true

	s.eventChannels.GameOverChan <- GameOverEvent{
		Board1:   s.board1,
		Board2:   s.board2,
		PlayerId: s.playerId,
		Winner:   message.Winner,
		Message:  message.Message,
	}
}

func (s *Session) Close() error {
	return s.connection.Close()
}
