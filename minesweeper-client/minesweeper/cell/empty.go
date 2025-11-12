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

func (c *EmptyCell) GetSnapshot() Snapshot {
	if c.IsOpened() {
		return OfEmpty()
	}
	if c.IsFlagged() {
		return OfFlag()
	}
	return OfUnchecked()
}

func (c *EmptyCell) ToggleFlag() {
	c.cellState.ToggleFlag()
}

func (c *EmptyCell) Open() {
	c.cellState.Open()
}
