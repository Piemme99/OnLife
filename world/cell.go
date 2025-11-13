package world

import "math/rand"

type CellType byte

type Cell struct {
	Type  CellType
	State int
}

const (
	Grass CellType = iota
	Fire
	Water
	Rock
	Life
)

const DefaultFireLifetime = 3

func NewCell(cellType CellType) Cell {
	cell := Cell{Type: cellType}
	if cellType == Fire {
		cell.State = DefaultFireLifetime
	}
	return cell
}

func RandomCellType() CellType {
	return CellType(rand.Intn(int(Life) + 1))
}
