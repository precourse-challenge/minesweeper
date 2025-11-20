package board

import (
	"fmt"
	"minesweeper-core/cell"
	"minesweeper-core/level"
	"minesweeper-core/position"
	"minesweeper-core/util"
)

type Board struct {
	cells         [][]cell.Cell
	landMineCount int
	flagCount     int
	status        Status
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
		status:        Ready,
	}
}

func (board *Board) InitializeGame() {
	board.status = InProgress
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
		board.status = Lose
		return nil
	}

	board.openSurroundedCells(cellPosition)

	if board.isGameWon() {
		board.status = Win
	}

	return nil
}

func (board *Board) IsInProgress() bool {
	return board.status == InProgress
}

func (board *Board) IsWinStatus() bool {
	return board.status == Win
}

func (board *Board) IsLoseStatus() bool {
	return board.status == Lose
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

func (board *Board) openSurroundedCells(cellPosition *position.CellPosition) {
	if board.isOpenedCell(cellPosition) {
		return
	}
	if board.isFlaggedCell(cellPosition) {
		return
	}
	if board.isLandMineCell(cellPosition) {
		return
	}

	board.openCell(cellPosition)

	if board.hasAdjacentLandMines(cellPosition) {
		return
	}

	surroundedPositions := board.findSurroundedPositions(cellPosition)
	for _, surroundedPosition := range surroundedPositions {
		board.openSurroundedCells(surroundedPosition)
	}
}

func (board *Board) findSurroundedPositions(cellPosition *position.CellPosition) []*position.CellPosition {
	surroundedPositions := make([]*position.CellPosition, 0, 8)

	for _, relativePosition := range position.RelativePositions {
		if cellPosition.CannotMoveBy(relativePosition) {
			continue
		}
		movedPosition := util.FatalIfError(cellPosition.MovedBy(relativePosition))
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

func (board *Board) openCell(cellPosition *position.CellPosition) {
	boardCell := board.cells[cellPosition.RowIndex()][cellPosition.ColIndex()]
	boardCell.Open()
}

func (board *Board) toggleFlagCell(cellPosition *position.CellPosition) {
	boardCell := board.cells[cellPosition.RowIndex()][cellPosition.ColIndex()]
	boardCell.ToggleFlag()
}

func (board *Board) isGameWon() bool {
	openedCellCount := board.countOpenedCells()

	fullSize := board.GetRowSize() * board.GetColSize()
	noLandMineCount := fullSize - board.landMineCount

	return openedCellCount == noLandMineCount
}

func (board *Board) countOpenedCells() int {
	var openedCellCount int
	for _, rowCells := range board.cells {
		for _, rowCell := range rowCells {
			if rowCell.IsOpened() {
				openedCellCount++
			}
		}
	}
	return openedCellCount
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
