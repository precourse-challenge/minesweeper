package cell

type EmptyCell struct {
}

func NewEmptyCell() *EmptyCell {
	return &EmptyCell{}
}

func (c *EmptyCell) IsLandMine() bool {
	return false
}
