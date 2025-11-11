package minesweeper

import (
	"minesweeper-client/minesweeper/level"
	"minesweeper-client/minesweeper/mode"
	"minesweeper-client/minesweeper/util"
	"minesweeper-client/minesweeper/view"
)

type Minesweeper struct {
}

func (minesweeper *Minesweeper) Run() {
	view.ShowGameStartMessage()

	gameLevel := util.Must(readInputGameLevel())
	gameMode := &mode.SingleMode{}

	gameMode.Start(gameLevel)
}

func readInputGameLevel() (level.GameLevel, error) {
	view.AskGameLevel()
	inputLevel := view.Read()
	view.ShowSelectedGameLevel(inputLevel)

	return level.From(inputLevel)
}
