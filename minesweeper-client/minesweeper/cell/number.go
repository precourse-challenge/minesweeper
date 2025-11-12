package cell

type NumberCell struct {
	cellState             *State
	adjacentLandMineCount int
}

func NewNumberCell(adjacentLandMineCount int) *NumberCell {
	return &NumberCell{
		cellState:             NewCellState(),
		adjacentLandMineCount: adjacentLandMineCount,
	}
}

func (c *NumberCell) IsLandMine() bool {
	return false
}

func (c *NumberCell) IsOpened() bool {
	return c.cellState.IsOpened()
}

func (c *NumberCell) IsFlagged() bool {
	return c.cellState.IsFlagged()
}

func (c *NumberCell) HasAdjacentLandMines() bool {
	return true
}

func (c *NumberCell) GetSnapshot() Snapshot {
	if c.IsOpened() {
		return OfNumber(c.adjacentLandMineCount)
	}
	if c.IsFlagged() {
		return OfFlag()
	}
	return OfUnchecked()
}
