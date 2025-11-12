package position

type RelativePosition struct {
	DeltaRow int
	DeltaCol int
}

var RelativePositions = []RelativePosition{
	{-1, -1},
	{-1, 0},
	{-1, 1},
	{0, -1},
	{0, 1},
	{1, -1},
	{1, 0},
	{1, 1},
}
