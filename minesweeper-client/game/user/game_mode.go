package user

import "fmt"

type GameMode string

const (
	Single          GameMode = "single"
	Multi           GameMode = "multi"
	UnknownGameMode GameMode = "unknown"
)

func GameModeFrom(input string) (GameMode, error) {
	switch input {
	case "single":
		return Single, nil
	case "multi":
		return Multi, nil
	default:
		return UnknownGameMode, fmt.Errorf("게임 모드는 (single / multi) 중 하나를 입력해야 합니다")
	}
}
