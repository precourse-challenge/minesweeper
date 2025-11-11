package minesweeper

import (
	"fmt"
	"minesweeper-client/minesweeper/level"
	"minesweeper-client/minesweeper/mode"
	"minesweeper-client/minesweeper/util"
	"minesweeper-client/minesweeper/view"
)

type Minesweeper struct {
}

func (minesweeper *Minesweeper) Run() {
	view.ShowGameStartMessage()

	gameLevel := util.Must(minesweeper.readInputGameLevel())
	gameMode := &mode.SingleMode{}

	gameMode.Start(gameLevel)
}

func (minesweeper *Minesweeper) readInputGameLevel() (level.GameLevel, error) {
	view.AskGameLevel()
	inputLevel := view.Read()
	view.ShowSelectedGameLevel(inputLevel)

	switch inputLevel {
	case "easy":
		return level.EasyLevel{}, nil
	case "normal":
		return level.NormalLevel{}, nil
	case "hard":
		return level.HardLevel{}, nil
	default:
		return nil, fmt.Errorf("난이도는 (easy / normal / hard) 중 하나를 입력해야 합니다")
	}
}
