package cell

type LandMineCell struct {
	cellState *State
}

func NewLandMineCell() *LandMineCell {
	return &LandMineCell{cellState: NewCellState()}
}

func (c *LandMineCell) IsLandMine() bool {
	return true
}

func (c *LandMineCell) IsOpened() bool {
	return c.cellState.IsOpened()
}

func (c *LandMineCell) IsFlagged() bool {
	return c.cellState.IsFlagged()
}

func (c *LandMineCell) HasAdjacentLandMines() bool {
	return false
}

func (c *LandMineCell) GetSnapshot() Snapshot {
	if c.IsOpened() {
		return OfLandMine()
	}
	if c.IsFlagged() {
		return OfFlag()
	}
	return OfUnchecked()
}
