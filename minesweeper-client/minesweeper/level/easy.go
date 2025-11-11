package level

type EasyLevel struct {
}

func (EasyLevel) Rows() int {
	return 9
}
func (EasyLevel) Cols() int {
	return 9
}
func (EasyLevel) MineCount() int {
	return 10
}
