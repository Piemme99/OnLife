// Package grid provides types and utilities to represent a 2D game grid.
package grid

import "fmt"

type CellType byte

type Cell struct {
	cellType CellType
}

const GridSize = 5

type Grid [GridSize][GridSize]CellType

const (
	Empty CellType = iota
	Grass
	Fire
	Water
	Rock
	Life
)

func (grid *Grid) Update() {
	var updatedGrid Grid
	for line, row := range grid {
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

func (grid *Grid) InitializeDefaultGrid() {
	*grid = Grid{
		{Water, Water, Water, Water, Water},
		{Water, Grass, Grass, Grass, Water},
		{Water, Grass, Fire, Grass, Water},
		{Water, Grass, Grass, Grass, Water},
		{Water, Water, Water, Water, Water},
	}
}

func (grid *Grid) Print() {
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

func (grid *Grid) hasFireNeigbour(x int, y int) bool {
	if grid[y+1][x] == Fire || grid[y-1][x] == Fire || grid[y][x+1] == Fire || grid[y][x-1] == Fire {
		return true
	}
	return false
}
