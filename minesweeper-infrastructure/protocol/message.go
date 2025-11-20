package protocol

import "minesweeper-infrastructure/dto"

type MessageType string

const (
	Join     MessageType = "join"
	Joined   MessageType = "joined"
	Start    MessageType = "start"
	GameOver MessageType = "game_over"
	Exit     MessageType = "exit"
	Error    MessageType = "error"
)

type Message struct {
	Type     MessageType  `json:"type"`
	PlayerId int          `json:"player_id"`
	Row      int          `json:"row,omitempty"`
	Col      int          `json:"col,omitempty"`
	Board1   dto.BoardDto `json:"board1,omitempty"`
	Board2   dto.BoardDto `json:"board2,omitempty"`
	Winner   int          `json:"winner,omitempty"`
	Message  string       `json:"message,omitempty"`
}
