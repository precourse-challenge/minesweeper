package minesweeper

import (
	"fmt"
	"minesweeper-client/minesweeper/level"
	"minesweeper-client/minesweeper/mode"
	"minesweeper-client/minesweeper/user"
	"minesweeper-client/minesweeper/view"
)

type Minesweeper struct {
}

func (minesweeper *Minesweeper) Run() {
	for {
		view.ShowGameStartMessage()

		gameLevel := readGameLevelWithRetry()

		gameMode := mode.NewSingleMode(gameLevel)
		gameMode.Start()

		userAction := minesweeper.readNextAction()

		if userAction == user.Retry {
			continue
		}
		if userAction == user.Exit {
			break
		}
	}
}

func readGameLevelWithRetry() level.GameLevel {
	for {
		selectedLevel, err := readInputGameLevel()
		if err != nil {
			view.ShowErrorMessage(err)
			continue
		}

		return selectedLevel
	}
}

func readInputGameLevel() (level.GameLevel, error) {
	view.AskGameLevel()
	inputLevel := view.Read()
	view.ShowSelectedGameLevel(inputLevel)

	return level.From(inputLevel)
}

func (minesweeper *Minesweeper) readNextAction() user.Action {
	for {
		view.ShowRestartMessage()

		inputAction := view.Read()
		action, err := minesweeper.parseUserAction(inputAction)
		if err != nil {
			view.ShowErrorMessage(err)
			continue
		}
		return action
	}
}

func (minesweeper *Minesweeper) parseUserAction(inputAction string) (user.Action, error) {
	userAction := user.From(inputAction)
	if userAction != user.Retry && userAction != user.Exit {
		return user.Unknown, fmt.Errorf("retry(재시작) 또는 exit(종료)을 입력해야 합니다")
	}
	return userAction, nil
}
