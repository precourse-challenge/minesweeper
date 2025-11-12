package minesweeper

import (
	"minesweeper-client/minesweeper/level"
	"minesweeper-client/minesweeper/mode"
	"minesweeper-client/minesweeper/view"
)

type Minesweeper struct {
}

func (minesweeper *Minesweeper) Run() {
	view.ShowGameStartMessage()

	gameLevel, err := readInputGameLevel()
	if err != nil {
		view.ShowErrorMessage(err)
	}

	gameMode := mode.NewSingleMode(gameLevel)
	gameMode.Start()
}

func readInputGameLevel() (level.GameLevel, error) {
	view.AskGameLevel()
	inputLevel := view.Read()
	view.ShowSelectedGameLevel(inputLevel)

	return level.From(inputLevel)
}
