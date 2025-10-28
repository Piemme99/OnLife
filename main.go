package main

import "fmt"

type Cell byte

const (
	Empty Cell = iota
	Grass
	Fire
	Water
)

const GridSize = 5

func main() {
	grid := defaultGrid()

	for range 5 {
		printGrid(grid)
		grid = updateGrid(grid)
	}
	printGrid(grid)
}

func updateGrid(grid [GridSize][GridSize]Cell) (updatedGrid [GridSize][GridSize]Cell) {
	for line, row := range grid {
		for column, c := range row {
			switch c {
			case Empty:
				updatedGrid[line][column] = c
			case Fire:
				updatedGrid[line][column] = c
			case Grass:
				if hasFireNeigbour(grid, line, column) {
					updatedGrid[line][column] = Fire
				} else {
					updatedGrid[line][column] = c
				}
			case Water:
				updatedGrid[line][column] = c
			}
		}
	}
	return updatedGrid
}

func hasFireNeigbour(grid [GridSize][GridSize]Cell, line int, column int) bool {
	if grid[line+1][column] == Fire || grid[line-1][column] == Fire || grid[line][column+1] == Fire || grid[line][column-1] == Fire {
		return true
	}
	return false
}

func defaultGrid() (grid [GridSize][GridSize]Cell) {
	return [GridSize][GridSize]Cell{
		{Water, Water, Water, Water, Water},
		{Water, Grass, Grass, Grass, Water},
		{Water, Grass, Fire, Grass, Water},
		{Water, Grass, Grass, Grass, Water},
		{Water, Water, Water, Water, Water},
	}
}

func printGrid(grid [GridSize][GridSize]Cell) {
	for _, row := range grid {
		for _, c := range row {
			switch c {
			case Empty:
				fmt.Print("⚪")
			case Fire:
				fmt.Print("🔥")
			case Grass:
				fmt.Print("🌱")
			case Water:
				fmt.Print("💧")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
