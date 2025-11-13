// Package grid provides types and utilities to represent a 2D game grid.
package world

import "fmt"

const GridSize = 5

type Grid [][]Cell

func NewGrid(width, height int) (grid Grid) {
	grid = make(Grid, height)

	for i := range grid {
		grid[i] = make([]Cell, width)
	}
	return grid
}

func (grid *Grid) Update() {
	if grid == nil {
		return
	}
	current := *grid
	if len(current) == 0 {
		return
	}
	height := len(current)
	width := len(current[0])
	updatedGrid := NewGrid(width, height)

	for y := range height {
		for x := range width {
			cell := current[y][x]
			switch cell.Type {
			case Rock, Life:
				lifeNeighbors := current.countLifeNeighbors(x, y)
				updatedGrid[y][x] = applyLifeRules(cell, lifeNeighbors)
			case Water:
				updatedGrid[y][x] = cell
			case Grass:
				if current.hasAdjacentType(x, y, Fire) {
					updatedGrid[y][x] = NewCell(Fire)
				} else {
					updatedGrid[y][x] = cell
				}
			case Fire:
				if current.hasAdjacentType(x, y, Water) {
					updatedGrid[y][x] = Cell{Type: Rock}
					continue
				}
				remaining := cell.State - 1
				if remaining <= 0 {
					updatedGrid[y][x] = Cell{Type: Rock}
				} else {
					updatedGrid[y][x] = Cell{Type: Fire, State: remaining}
				}
			default:
				updatedGrid[y][x] = cell
			}
		}
	}
	*grid = updatedGrid
}

func (grid Grid) InitializeRandomGrid() {
	for line, row := range grid {
		for column := range row {
			grid[line][column] = NewCell(RandomCellType())
		}
	}
}

func (grid Grid) Print() {
	for _, row := range grid {
		for _, cell := range row {
			switch cell.Type {
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

func (grid Grid) hasAdjacentType(x, y int, target CellType) bool {
	found := false
	grid.forEachNeighbor(x, y, func(nx, ny int) {
		if grid[ny][nx].Type == target {
			found = true
		}
	})
	return found
}

func applyLifeRules(cell Cell, lifeNeighbors int) Cell {
	switch cell.Type {
	case Life:
		if lifeNeighbors == 2 || lifeNeighbors == 3 {
			return cell
		}
		return Cell{Type: Rock}
	case Rock:
		if lifeNeighbors == 3 {
			return Cell{Type: Life}
		}
	}
	return cell
}

func (grid Grid) countLifeNeighbors(x, y int) int {
	if len(grid) == 0 {
		return 0
	}
	count := 0
	directions := [8][2]int{
		{-1, -1},
		{0, -1},
		{1, -1},
		{-1, 0},
		{1, 0},
		{-1, 1},
		{0, 1},
		{1, 1},
	}
	for _, dir := range directions {
		nx := x + dir[0]
		ny := y + dir[1]
		if ny < 0 || ny >= len(grid) {
			continue
		}
		if nx < 0 || nx >= len(grid[ny]) {
			continue
		}
		if grid[ny][nx].Type == Life {
			count++
		}
	}
	return count
}

func (grid Grid) forEachNeighbor(x, y int, fn func(nx, ny int)) {
	if len(grid) == 0 {
		return
	}
	directions := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	for _, dir := range directions {
		nx := x + dir[0]
		ny := y + dir[1]
		if ny < 0 || ny >= len(grid) {
			continue
		}
		if nx < 0 || nx >= len(grid[ny]) {
			continue
		}
		fn(nx, ny)
	}
}
