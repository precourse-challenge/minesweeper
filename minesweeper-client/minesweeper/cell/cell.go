package cell

type Cell interface {
	IsLandMine() bool
	IsOpened() bool
	IsFlagged() bool
	HasAdjacentLandMines() bool
	GetSnapshot() Snapshot
}
