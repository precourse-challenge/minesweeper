package position

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_입력받은_좌표로_cellPosition을_생성한다(t *testing.T) {
	// given
	inputs := []struct {
		row int
		col int
	}{
		{0, 0},
		{2, 3},
		{5, 10},
	}

	// when & then
	for _, input := range inputs {
		cellPosition, _ := NewCellPosition(input.row, input.col)
		assert.NotNil(t, cellPosition)
	}
}

func Test_좌표가_0보다_작은_값이_존재하면_예외가_발생한다(t *testing.T) {
	// given
	inputs := []struct {
		row int
		col int
	}{
		{-1, 0},
		{0, -5},
		{-3, -2},
	}

	// when & then
	for _, input := range inputs {
		_, err := NewCellPosition(input.row, input.col)
		assert.EqualError(t, err, "올바르지 않은 좌표값입니다")
	}
}
