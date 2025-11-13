package main

import (
	"log"
	"os"

	"OnLife/world"
)

const defaultIterations = 10

func main() {
	args := os.Args[1:]
	var (
		scenarioPath string
		savePath     string
	)
	if len(args) > 0 {
		scenarioPath = args[0]
	}
	if len(args) > 1 {
		savePath = args[1]
	}

	var grid world.Grid

	if scenarioPath != "" {
		scenario, err := world.LoadScenario(scenarioPath)
		if err != nil {
			log.Fatalf("load scenario: %v", err)
		}
		loadedGrid, err := scenario.BuildGrid()
		if err != nil {
			log.Fatalf("build scenario grid: %v", err)
		}
		if len(loadedGrid) == 0 || len(loadedGrid[0]) == 0 {
			log.Fatal("scenario grid is empty")
		}
		grid = loadedGrid
	} else {
		w := world.NewWord(world.GridSize, world.GridSize)
		grid = w.Grid
		grid.InitializeRandomGrid()
	}

	for range defaultIterations {
		grid.Print()
		grid.Update()
	}
	grid.Print()

	if savePath != "" {
		snapshot := world.GridToScenario("snapshot", grid)
		if err := world.SaveScenario(savePath, snapshot); err != nil {
			log.Fatalf("save scenario: %v", err)
		}
	}
}
