package mode

import (
	"minesweeper-client/minesweeper/board"
	"minesweeper-client/minesweeper/level"
)

type SingleMode struct {
	board *board.Board
}

func NewSingleMode(level level.GameLevel) *SingleMode {
	return &SingleMode{board: board.NewBoard(level)}
}

func (mode *SingleMode) Start() {
	mode.board.InitializeGame()
}
