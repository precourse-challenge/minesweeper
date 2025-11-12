package board

import (
	"fmt"
	"minesweeper-client/minesweeper/cell"
	"minesweeper-client/minesweeper/level"
	"minesweeper-client/minesweeper/position"
	"minesweeper-client/minesweeper/util"
)

type Board struct {
	cells         [][]cell.Cell
	landMineCount int
	flagCount     int
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

func (board *Board) Flag(cellPosition *position.CellPosition) error {
	if board.isOpenedCell(cellPosition) {
		return fmt.Errorf("열려있는 셀에 깃발을 꽂을 수 없습니다")
	}

	if board.isFlaggedCell(cellPosition) {
		board.toggleFlagCell(cellPosition)
		board.flagCount--
		return nil
	}

	if board.hasNoRemainingFlags() {
		return fmt.Errorf("더 이상 깃발을 꽂을 수 없습니다")
	}

	board.toggleFlagCell(cellPosition)
	board.flagCount++
	return nil
}

func (board *Board) Open(cellPosition *position.CellPosition) error {
	if board.isOpenedCell(cellPosition) {
		return fmt.Errorf("이미 열려있는 셀 입니다")
	}

	if board.isFlaggedCell(cellPosition) {
		return fmt.Errorf("셀을 열려면 깃발을 먼저 해제해야 합니다")
	}

	if board.isLandMineCell(cellPosition) {
		board.openCell(cellPosition)
		board.gameStatus = Lose
		return nil
	}

	board.openCell(cellPosition)
	return nil
}

func (board *Board) openCell(cellPosition *position.CellPosition) {
	boardCell := board.cells[cellPosition.RowIndex()][cellPosition.ColIndex()]
	boardCell.Open()
}

func (board *Board) toggleFlagCell(cellPosition *position.CellPosition) {
	boardCell := board.cells[cellPosition.RowIndex()][cellPosition.ColIndex()]
	boardCell.ToggleFlag()
}

func (board *Board) IsInProgress() bool {
	return board.gameStatus == InProgress
}

func (board *Board) IsOutOfBounds(cellPosition *position.CellPosition) bool {
	rowSize := board.GetRowSize()
	colSize := board.GetColSize()

	return cellPosition.RowIndex() < 0 || cellPosition.ColIndex() < 0 ||
		cellPosition.RowIndex() >= rowSize || cellPosition.ColIndex() >= colSize
}

func (board *Board) GetRemainingFlags() int {
	return board.landMineCount - board.flagCount
}

func (board *Board) GetSnapshot(cellPosition *position.CellPosition) cell.Snapshot {
	boardCell := board.cells[cellPosition.RowIndex()][cellPosition.ColIndex()]
	return boardCell.GetSnapshot()
}

func (board *Board) GetColSize() int {
	return len(board.cells[0])
}

func (board *Board) GetRowSize() int {
	return len(board.cells)
}

func (board *Board) initializeEmptyCells() {
	for rowIndex, row := range board.cells {
		for colIndex := range row {
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
		movedPosition := util.FatalIfError(cellPosition.MovedBy(surroundedPosition))
		if board.IsOutOfBounds(movedPosition) {
			continue
		}
		surroundedPositions = append(surroundedPositions, movedPosition)
	}
	return surroundedPositions
}

func (board *Board) updateCell(cellPosition *position.CellPosition, cell cell.Cell) {
	board.cells[cellPosition.RowIndex()][cellPosition.ColIndex()] = cell
}

func (board *Board) isOpenedCell(cellPosition *position.CellPosition) bool {
	boardCell := board.cells[cellPosition.RowIndex()][cellPosition.ColIndex()]
	return boardCell.IsOpened()
}

func (board *Board) isLandMineCell(cellPosition *position.CellPosition) bool {
	boardCell := board.cells[cellPosition.RowIndex()][cellPosition.ColIndex()]
	return boardCell.IsLandMine()
}

func (board *Board) isFlaggedCell(cellPosition *position.CellPosition) bool {
	boardCell := board.cells[cellPosition.RowIndex()][cellPosition.ColIndex()]
	return boardCell.IsFlagged()
}

func (board *Board) hasAdjacentLandMines(cellPosition *position.CellPosition) bool {
	boardCell := board.cells[cellPosition.RowIndex()][cellPosition.ColIndex()]
	return boardCell.HasAdjacentLandMines()
}

func (board *Board) hasNoRemainingFlags() bool {
	return board.flagCount >= board.landMineCount
}
