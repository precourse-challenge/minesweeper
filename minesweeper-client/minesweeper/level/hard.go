package level

type HardLevel struct {
}

func (HardLevel) Rows() int {
	return 20
}
func (HardLevel) Cols() int {
	return 20
}
func (HardLevel) MineCount() int {
	return 80
}
