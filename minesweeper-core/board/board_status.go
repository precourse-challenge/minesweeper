package board

type Status int

const (
	Ready Status = iota
	InProgress
	Win
	Lose
)
