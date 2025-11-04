package world

import "math/rand"

type CellType byte

type Cell struct {
	cellType CellType
}

const (
	Empty CellType = iota
	Grass
	Fire
	Water
	Rock
	Life
)

func RandomCell() (cell CellType) {
	return CellType(rand.Intn(int(Life) + 1))
}
