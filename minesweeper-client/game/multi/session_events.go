package multi

import "minesweeper-infrastructure/dto"

type JoinedEvent struct {
	PlayerId int
}

type StartEvent struct {
	Board1   dto.BoardDto
	Board2   dto.BoardDto
	PlayerId int
}

type UpdateEvent struct {
	Board1   dto.BoardDto
	Board2   dto.BoardDto
	PlayerId int
}

type ErrorEvent struct {
	Err error
}

type GameOverEvent struct {
	Board1   dto.BoardDto
	Board2   dto.BoardDto
	PlayerId int
	Winner   int
}
