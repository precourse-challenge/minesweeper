package match

import (
	"fmt"
	"minesweeper-core/board"
	"minesweeper-core/cell"
	"minesweeper-core/level"
	"minesweeper-core/position"
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

func (m *Match) Open(playerId int, position *position.CellPosition) (Result, error) {
	if m.status != Playing {
		return Result{IsGameOver: false}, fmt.Errorf("매치가 진행상태가 아닙니다")
	}

	playerBoard := m.getPlayerBoard(playerId)
	err := playerBoard.Open(position)
	if err != nil {
		return Result{IsGameOver: false}, err
	}

	if playerBoard.IsLoseStatus() {
		m.status = Finished
		m.winner = m.getOpponentId(playerId)
		return Result{
			IsGameOver: true,
			Winner:     m.winner,
			Message:    "지뢰를 밟았습니다.",
		}, nil
	}

	if playerBoard.IsWinStatus() {
		m.status = Finished
		m.winner = playerId
		return Result{
			IsGameOver: true,
			Winner:     m.winner,
			Message:    "게임이 끝났습니다.",
		}, nil
	}

	return Result{IsGameOver: false}, nil
}

func (m *Match) Flag(playerId int, position *position.CellPosition) error {
	if m.status != Playing {
		return fmt.Errorf("매치가 진행상태가 아닙니다")
	}
	playerBoard := m.getPlayerBoard(playerId)
	return playerBoard.Flag(position)
}

func (m *Match) GetPlayer1Board() [][]cell.Snapshot {
	return m.player1Board.GetSnapshots()
}

func (m *Match) GetPlayer2Board() [][]cell.Snapshot {
	return m.player2Board.GetSnapshots()
}

func (m *Match) getPlayerBoard(playerId int) *board.Board {
	if playerId == 1 {
		return m.player1Board
	}
	return m.player2Board
}

func (m *Match) getOpponentId(playerId int) int {
	if playerId == 1 {
		return 2
	}
	return 1
}
