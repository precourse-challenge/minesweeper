package multi

import "minesweeper-infrastructure/dto"

type JoinedEvent struct {
	PlayerId int
	Message  string
}

type StartEvent struct {
	Board1   dto.BoardDto
	Board2   dto.BoardDto
	PlayerId int
	Message  string
}

type UpdateEvent struct {
	Board1   dto.BoardDto
	Board2   dto.BoardDto
	PlayerId int
	Message  string
}

type ErrorEvent struct {
	Err     error
	Message string
}

type GameOverEvent struct {
	Board1   dto.BoardDto
	Board2   dto.BoardDto
	PlayerId int
	Winner   int
	Message  string
}
