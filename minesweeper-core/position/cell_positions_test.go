package position

import (
	"minesweeper-core/cell"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_요청_개수만큼_무작위의_셀을_반환한다(t *testing.T) {
	// given
	board := [][]cell.Cell{
		{cell.NewEmptyCell(), cell.NewEmptyCell()},
		{cell.NewEmptyCell(), cell.NewEmptyCell()},
	}
	positions := NewCellPositions(board)

	// when
	extractedPositions := positions.ExtractRandomPositions(1)

	// then
	assert.Len(t, extractedPositions, 1)
}

func Test_요청_개수가_전체_셀_개수를_초과하면_전체_셀의_개수만큼_반환한다(t *testing.T) {
	// given
	board := [][]cell.Cell{
		{cell.NewEmptyCell(), cell.NewEmptyCell()},
		{cell.NewEmptyCell(), cell.NewEmptyCell()},
	}
	positions := NewCellPositions(board)
	positionCount := len(positions.positions)

	// when
	extractedPositions := positions.ExtractRandomPositions(5)

	// then
	assert.Len(t, extractedPositions, positionCount)
}

func Test_특정_셀을_제외한_나머지_셀들을_반환한다(t *testing.T) {
	// given
	board := [][]cell.Cell{
		{cell.NewEmptyCell(), cell.NewEmptyCell()},
		{cell.NewEmptyCell(), cell.NewEmptyCell()},
	}
	positions := NewCellPositions(board)
	toSubtract := []*CellPosition{mustCellPosition(t, 0, 1)}

	// when
	subtractedPositions := positions.Subtract(toSubtract)

	// then
	assert.True(t, contains(subtractedPositions, mustCellPosition(t, 0, 0)))
	assert.True(t, contains(subtractedPositions, mustCellPosition(t, 1, 0)))
	assert.True(t, contains(subtractedPositions, mustCellPosition(t, 1, 1)))
	assert.False(t, contains(subtractedPositions, mustCellPosition(t, 0, 1)))
}

func mustCellPosition(t *testing.T, row, col int) *CellPosition {
	t.Helper()
	position, err := NewCellPosition(row, col)
	if err != nil {
		t.Fatal(err)
	}
	return position
}
