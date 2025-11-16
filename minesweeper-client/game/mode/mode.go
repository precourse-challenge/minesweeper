package mode

import "minesweeper-core/level"

type GameMode interface {
	Start(gameLevel level.GameLevel)
}
