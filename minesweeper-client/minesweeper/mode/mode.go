package mode

import "minesweeper-client/minesweeper/level"

type GameMode interface {
	Start(gameLevel level.GameLevel)
}
