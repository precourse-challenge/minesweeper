package level

type GameLevel interface {
	Rows() int
	Cols() int
	MineCount() int
}
