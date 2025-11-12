package position

type RelativePosition struct {
	DeltaRow int
	DeltaCol int
}

var SurroundedPositions = []RelativePosition{
	{-1, -1},
	{-1, 0},
	{-1, 1},
	{0, -1},
	{0, 1},
	{1, -1},
	{1, 0},
	{1, 1},
}
