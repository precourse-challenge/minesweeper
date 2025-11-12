package cell

type Snapshot struct {
	status                SnapshotStatus
	adjacentLandMineCount int
}

func NewCellSnapshot(status SnapshotStatus, adjacentLandMineCount int) Snapshot {
	return Snapshot{status: status, adjacentLandMineCount: adjacentLandMineCount}
}

func OfUnchecked() Snapshot {
	return NewCellSnapshot(Unchecked, 0)
}

func OfEmpty() Snapshot {
	return NewCellSnapshot(Empty, 0)
}

func OfFlag() Snapshot {
	return NewCellSnapshot(Flag, 0)
}

func OfLandMine() Snapshot {
	return NewCellSnapshot(LandMine, 0)
}

func OfNumber(adjacentLandMineCount int) Snapshot {
	return NewCellSnapshot(Number, adjacentLandMineCount)
}

func (snapshot Snapshot) GetStatus() SnapshotStatus {
	return snapshot.status
}

func (snapshot Snapshot) GetAdjacentLandMineCount() int {
	return snapshot.adjacentLandMineCount
}
