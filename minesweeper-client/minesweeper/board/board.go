package board

import (
	"minesweeper-client/minesweeper/cell"
	"minesweeper-client/minesweeper/level"
	"minesweeper-client/minesweeper/position"
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

func (board *Board) InitializeGame() {
	board.gameStatus = InProgress
	board.initializeEmptyCells()

	cellPositions := position.NewCellPositions(board.cells)
	landMinePositions := cellPositions.ExtractRandomPositions(board.landMineCount)
	board.initializeLandMineCells(landMinePositions)
}

func (board *Board) initializeEmptyCells() {
	for rowIndex, row := range board.cells {
		for colIndex, _ := range row {
			board.cells[rowIndex][colIndex] = cell.NewEmptyCell()
		}
	}
}

func (board *Board) initializeLandMineCells(cellPositions []*position.CellPosition) {
	for _, cellPosition := range cellPositions {
		board.updateCell(cellPosition, cell.NewLandMineCell())
	}
}

func (board *Board) updateCell(cellPosition *position.CellPosition, cell cell.Cell) {
	board.cells[cellPosition.RowIndex()][cellPosition.ColIndex()] = cell
}
