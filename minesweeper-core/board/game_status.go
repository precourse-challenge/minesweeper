package board

type GameStatus int

const (
	Ready GameStatus = iota
	InProgress
	Win
	Lose
)
