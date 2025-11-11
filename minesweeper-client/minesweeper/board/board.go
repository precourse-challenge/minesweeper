package board

import (
	"minesweeper-client/minesweeper/cell"
	"minesweeper-client/minesweeper/level"
)

type Board struct {
	cells         [][]cell.Cell
	landMineCount int
	gameStatus    GameStatus
}

func NewBoard(level level.GameLevel) *Board {
	rows := level.Rows()
	cols := level.Cols()

	cells := make([][]cell.Cell, rows)
	for i := 0; i < rows; i++ {
		cells[i] = make([]cell.Cell, cols)
	}

	return &Board{
		cells:         cells,
		landMineCount: level.MineCount(),
		gameStatus:    Ready,
	}
}
