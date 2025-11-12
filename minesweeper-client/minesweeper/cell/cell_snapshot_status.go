package cell

type SnapshotStatus int

const (
	Unchecked SnapshotStatus = iota
	Empty
	Flag
	LandMine
	Number
)
