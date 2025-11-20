package position

import (
	"math/rand"
	"minesweeper-core/cell"
	"minesweeper-core/util"
	"time"
)

type CellPositions struct {
	positions []*CellPosition
}

func NewCellPositions(board [][]cell.Cell) *CellPositions {
	var positions []*CellPosition

	for rowIndex := range board {
		for colIndex := range board[rowIndex] {
			position := util.FatalIfError(NewCellPosition(rowIndex, colIndex))
			positions = append(positions, position)
		}
	}
	return &CellPositions{positions: positions}
}

func (positions *CellPositions) ExtractRandomPositions(count int) []*CellPosition {
	cellPositions := make([]*CellPosition, len(positions.positions))
	copy(cellPositions, positions.positions)

	positions.shufflePositions(cellPositions)

	if count > len(cellPositions) {
		count = len(cellPositions)
	}
	return cellPositions[:count]
}

func (positions *CellPositions) Subtract(toSubtract []*CellPosition) []*CellPosition {
	cellPositions := make([]*CellPosition, 0, len(positions.positions))

	for _, position := range positions.positions {
		if doesNotContain(toSubtract, position) {
			cellPositions = append(cellPositions, position)
		}
	}
	return cellPositions
}

func (positions *CellPositions) shufflePositions(cellPositions []*CellPosition) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	random.Shuffle(len(cellPositions), func(a, b int) {
		cellPositions[a], cellPositions[b] = cellPositions[b], cellPositions[a]
	})
}

func doesNotContain(positions []*CellPosition, position *CellPosition) bool {
	return !contains(positions, position)
}

func contains(cellPositions []*CellPosition, position *CellPosition) bool {
	for _, cellPosition := range cellPositions {
		if cellPosition.Equals(position) {
			return true
		}
	}
	return false
}
