package util

import (
	"minesweeper-client/minesweeper/view"
	"os"
)

func Must[T any](value T, err error) T {
	if err != nil {
		view.ShowErrorMessage(err)
		os.Exit(1)
	}
	return value
}
