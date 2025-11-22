package dto

import (
	"minesweeper-core/cell"
)

type CellSnapshotDto struct {
	Status cell.SnapshotStatus `json:"status"`
	Number int                 `json:"number"`
}

func ToCellSnapshotDto(snapshot cell.Snapshot) CellSnapshotDto {
	status := snapshot.GetStatus()
	number := snapshot.GetAdjacentLandMineCount()

	if status == cell.LandMine && snapshot.IsUnchecked() {
		status = cell.Unchecked
		number = 0
	}

	return CellSnapshotDto{
		Status: status,
		Number: number,
	}
}
