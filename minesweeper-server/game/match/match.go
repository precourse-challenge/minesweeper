package match

import (
	"minesweeper-core/board"
	"minesweeper-core/cell"
	"minesweeper-core/level"
)

type Match struct {
	player1Board *board.Board
	player2Board *board.Board
	status       Status
	winner       int
}

func NewMatch() *Match {
	board1 := board.NewBoard(level.EasyLevel{})
	board2 := board.NewBoard(level.EasyLevel{})
	return &Match{
		player1Board: board1,
		player2Board: board2,
		status:       Ready,
		winner:       0,
	}
}

func (m *Match) InitializeGame() {
	m.player1Board.InitializeGame()
	m.player2Board.InitializeGame()
	m.status = Playing
}

func (m *Match) GetPlayer1Board() [][]cell.Snapshot {
	return m.player1Board.GetSnapshots()
}

func (m *Match) GetPlayer2Board() [][]cell.Snapshot {
	return m.player2Board.GetSnapshots()
}
