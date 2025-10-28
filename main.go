package main

import "OnLife/grid"

func main() {
	var grid grid.Grid

	grid.InitializeDefaultGrid()
	for range 5 {
		grid.Print()
		grid.Update()
	}
	grid.Print()
}
