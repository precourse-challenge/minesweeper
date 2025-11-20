package dto

import (
	"minesweeper-core/cell"
)

type BoardDto [][]CellSnapshotDto

func ToBoardDto(snapshots [][]cell.Snapshot) BoardDto {
	rowSize := len(snapshots)
	colSize := len(snapshots[0])

	boardDto := make(BoardDto, rowSize)
	for row := 0; row < rowSize; row++ {
		boardDto[row] = make([]CellSnapshotDto, colSize)
		for col := 0; col < colSize; col++ {
			boardDto[row][col] = ToCellSnapshotDto(snapshots[row][col])
		}
	}
	return boardDto
}
