package level

import "fmt"

type GameLevel interface {
	Rows() int
	Cols() int
	MineCount() int
}

func From(inputLevel string) (GameLevel, error) {
	switch inputLevel {
	case "easy":
		return EasyLevel{}, nil
	case "normal":
		return NormalLevel{}, nil
	case "hard":
		return HardLevel{}, nil
	default:
		return nil, fmt.Errorf("난이도는 (easy / normal / hard) 중 하나를 입력해야 합니다")
	}
}
