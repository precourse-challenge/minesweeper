package game

import (
	"fmt"
	"minesweeper-client/game/multi"
	"minesweeper-client/game/single"
	"minesweeper-client/game/user"
	"minesweeper-client/game/view"
	"minesweeper-core/level"
)

func StartGameClient() {
	view.ShowGameStartMessage()

	mode := readGameMode()

	for {
		switch mode {
		case user.Single:
			runSingleMode()
		case user.Multi:
			runMultiMode()
		}

		userAction := readNextAction()

		if userAction == user.Retry {
			continue
		}
		if userAction == user.Exit {
			break
		}
	}
}

func readGameMode() user.GameMode {
	for {
		gameMode, err := readInputGameMode()
		if err != nil {
			view.ShowErrorMessage(err)
			continue
		}

		return gameMode
	}
}

func readInputGameMode() (user.GameMode, error) {
	view.ShowGameModeSelection()
	inputGameMode := view.Read()
	return user.GameModeFrom(inputGameMode)
}

func runSingleMode() {
	gameLevel := readGameLevel()
	singleMode := single.NewSingleMode(gameLevel)
	singleMode.Start()
}

func runMultiMode() {
	multiMode := multi.NewMultiMode()
	multiMode.Start()
}

func readGameLevel() level.GameLevel {
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

func readNextAction() user.Action {
	for {
		view.ShowRestartMessage()

		inputAction := view.Read()
		action, err := parseUserAction(inputAction)
		if err != nil {
			view.ShowErrorMessage(err)
			continue
		}
		return action
	}
}

func parseUserAction(inputAction string) (user.Action, error) {
	userAction := user.ActionFrom(inputAction)
	if userAction != user.Retry && userAction != user.Exit {
		return user.UnknownAction, fmt.Errorf("retry(재시작) 또는 exit(종료)을 입력해야 합니다")
	}
	return userAction, nil
}
