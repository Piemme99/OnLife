// Package grid provides types and utilities to represent a 2D game grid.
package world

import "fmt"

const GridSize = 5

// TODO: Grid should be 1D for performances
type Grid [][]CellType

func NewGrid(width, height int) (grid Grid) {
	grid = make(Grid, height)

	for i := range grid {
		grid[i] = make([]CellType, width)
	}
	return grid
}

// TODO:  Redo Update
func (grid *Grid) Update() {
	updatedGrid := NewGrid(GridSize, GridSize)
	for line, row := range *grid {
		for column, cell := range row {
			switch cell {
			case Empty:
				updatedGrid[line][column] = cell
			case Fire:
				updatedGrid[line][column] = cell
			case Grass:
				if grid.hasFireNeigbour(line, column) {
					updatedGrid[line][column] = Fire
				} else {
					updatedGrid[line][column] = cell
				}
			case Water:
				updatedGrid[line][column] = cell
			}
		}
	}
	*grid = updatedGrid
}

func (grid Grid) InitializeRandomGrid() {
	for line, row := range grid {
		for column := range row {
			grid[line][column] = RandomCell()
		}
	}
}

func (grid Grid) Print() {
	for _, row := range grid {
		for _, cell := range row {
			switch cell {
			case Empty:
				fmt.Print("⚪")
			case Fire:
				fmt.Print("🔥")
			case Grass:
				fmt.Print("🌱")
			case Water:
				fmt.Print("💧")
			case Rock:
				fmt.Print("🪨")
			case Life:
				fmt.Print("🧡")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func (grid Grid) hasFireNeigbour(x int, y int) bool {
	if grid[y+1][x] == Fire || grid[y-1][x] == Fire || grid[y][x+1] == Fire || grid[y][x-1] == Fire {
		return true
	}
	return false
}
