package level

type NormalLevel struct {
}

func (NormalLevel) Rows() int {
	return 15
}
func (NormalLevel) Cols() int {
	return 15
}
func (NormalLevel) MineCount() int {
	return 40
}
