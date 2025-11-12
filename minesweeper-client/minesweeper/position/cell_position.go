package position

import "fmt"

type CellPosition struct {
	rowIndex int
	colIndex int
}

func NewCellPosition(rowIndex, colIndex int) (*CellPosition, error) {
	if rowIndex < 0 || colIndex < 0 {
		return nil, fmt.Errorf("올바르지 않은 좌표값입니다")
	}
	return &CellPosition{rowIndex: rowIndex, colIndex: colIndex}, nil
}

func (position *CellPosition) CannotMoveBy(relativePosition RelativePosition) bool {
	movedRow := position.rowIndex + relativePosition.DeltaRow
	movedCol := position.colIndex + relativePosition.DeltaCol

	return movedRow < 0 || movedCol < 0
}

func (position *CellPosition) MovedBy(relativePosition RelativePosition) (*CellPosition, error) {
	movedRow := position.rowIndex + relativePosition.DeltaRow
	movedCol := position.colIndex + relativePosition.DeltaCol

	return NewCellPosition(movedRow, movedCol)
}

func (position *CellPosition) Equals(otherPosition *CellPosition) bool {
	return position.rowIndex == otherPosition.rowIndex &&
		position.colIndex == otherPosition.colIndex
}

func (position *CellPosition) RowIndex() int {
	return position.rowIndex
}

func (position *CellPosition) ColIndex() int {
	return position.colIndex
}
