package main

import "OnLife/world"

func main() {
	world := world.NewWord(world.GridSize, world.GridSize)
	grid := world.Grid

	grid.InitializeRandomGrid()
	for range 5 {
		grid.Print()
		grid.Update()
	}
	grid.Print()
}
