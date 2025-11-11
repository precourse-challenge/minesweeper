package cell

type LandMineCell struct {
}

func NewLandMineCell() *LandMineCell {
	return &LandMineCell{}
}

func (c *LandMineCell) IsLandMine() bool {
	return true
}
