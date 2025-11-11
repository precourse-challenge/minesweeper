package board

import (
	"minesweeper-client/minesweeper/cell"
	"minesweeper-client/minesweeper/level"
	"minesweeper-client/minesweeper/position"
	"minesweeper-client/minesweeper/util"
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

	numberPositions := cellPositions.Subtract(landMinePositions)
	board.initializeNumberCells(numberPositions)
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

func (board *Board) initializeNumberCells(cellPositions []*position.CellPosition) {
	for _, cellPosition := range cellPositions {
		adjacentLandMineCount := board.countAdjacentLandMines(cellPosition)
		if adjacentLandMineCount > 0 {
			board.updateCell(cellPosition, cell.NewNumberCell(adjacentLandMineCount))
		}
	}
}

func (board *Board) countAdjacentLandMines(cellPosition *position.CellPosition) int {
	surroundedPositions := board.findSurroundedPositions(cellPosition)

	adjacentLandMineCount := 0
	for _, surroundedPosition := range surroundedPositions {
		if board.isLandMineCell(surroundedPosition) {
			adjacentLandMineCount++
		}
	}
	return adjacentLandMineCount
}

func (board *Board) findSurroundedPositions(cellPosition *position.CellPosition) []*position.CellPosition {
	var surroundedPositions []*position.CellPosition

	for _, surroundedPosition := range position.SurroundedPositions {
		if cellPosition.CannotMoveBy(surroundedPosition) {
			continue
		}
		movedPosition := util.Must(cellPosition.MovedBy(surroundedPosition))
		if board.isOutOfBounds(movedPosition) {
			continue
		}
		surroundedPositions = append(surroundedPositions, movedPosition)
	}
	return surroundedPositions
}

func (board *Board) updateCell(cellPosition *position.CellPosition, cell cell.Cell) {
	board.cells[cellPosition.RowIndex()][cellPosition.ColIndex()] = cell
}

func (board *Board) isLandMineCell(cellPosition *position.CellPosition) bool {
	boardCell := board.cells[cellPosition.RowIndex()][cellPosition.ColIndex()]
	return boardCell.IsLandMine()
}

func (board *Board) isOutOfBounds(cellPosition *position.CellPosition) bool {
	rowSize := len(board.cells)
	colSize := len(board.cells[0])

	return cellPosition.RowIndex() < 0 || cellPosition.ColIndex() < 0 ||
		cellPosition.RowIndex() >= rowSize || cellPosition.ColIndex() >= colSize
}
