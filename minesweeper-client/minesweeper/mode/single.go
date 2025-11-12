package mode

import (
	"minesweeper-client/minesweeper/board"
	"minesweeper-client/minesweeper/level"
	"minesweeper-client/minesweeper/view"
)

type SingleMode struct {
	board *board.Board
}

func NewSingleMode(level level.GameLevel) *SingleMode {
	return &SingleMode{board: board.NewBoard(level)}
}

func (m *SingleMode) Start() {
	m.board.InitializeGame()
	view.ShowBoard(m.board)
}
