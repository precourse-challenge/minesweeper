package position

import (
	"math/rand"
	"minesweeper-client/minesweeper/cell"
	"minesweeper-client/minesweeper/util"
	"time"
)

type CellPositions struct {
	positions []*CellPosition
}

func NewCellPositions(board [][]cell.Cell) *CellPositions {
	var positions []*CellPosition

	for rowIndex := range board {
		for colIndex := range board[rowIndex] {
			position := util.Must(NewCellPosition(rowIndex, colIndex))
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

func (positions *CellPositions) shufflePositions(cellPositions []*CellPosition) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	random.Shuffle(len(cellPositions), func(a, b int) {
		cellPositions[a], cellPositions[b] = cellPositions[b], cellPositions[a]
	})
}
