package world

type World struct {
	width, height int
	Grid          Grid
}

func NewWord(width, height int) (world World) {
	world.height = height
	world.width = width
	world.Grid = NewGrid(width, height)
	return world
}
