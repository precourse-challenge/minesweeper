package cell

type State struct {
	isOpened  bool
	isFlagged bool
}

func NewCellState() *State {
	return &State{
		isOpened:  false,
		isFlagged: false,
	}
}

func (c *State) IsOpened() bool {
	return c.isOpened
}

func (c *State) IsFlagged() bool {
	return c.isFlagged
}
