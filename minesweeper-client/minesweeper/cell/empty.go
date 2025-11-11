package cell

type EmptyCell struct {
	cellState *State
}

func NewEmptyCell() *EmptyCell {
	return &EmptyCell{cellState: NewCellState()}
}

func (c *EmptyCell) IsLandMine() bool {
	return false
}

func (c *EmptyCell) IsOpened() bool {
	return c.cellState.IsOpened()
}

func (c *EmptyCell) IsFlagged() bool {
	return c.cellState.IsFlagged()
}

func (c *EmptyCell) HasAdjacentLandMines() bool {
	return false
}
