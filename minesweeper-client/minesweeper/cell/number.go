package cell

type NumberCell struct {
	adjacentLandMineCount int
}

func NewNumberCell(adjacentLandMineCount int) *NumberCell {
	return &NumberCell{
		adjacentLandMineCount: adjacentLandMineCount,
	}
}

func (c *NumberCell) IsLandMine() bool {
	return false
}
