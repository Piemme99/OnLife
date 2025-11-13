package world

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Scenario represents a serialized grid layout along with optional metadata.
type Scenario struct {
	Name         string   `json:"name,omitempty"`
	Description  string   `json:"description,omitempty"`
	FireLifetime int      `json:"fireLifetime,omitempty"`
	Rows         []string `json:"rows"`
}

// Width returns the number of columns described by the scenario.
func (s Scenario) Width() int {
	if len(s.Rows) == 0 {
		return 0
	}
	return len([]rune(s.Rows[0]))
}

// Height returns the number of rows described by the scenario.
func (s Scenario) Height() int {
	return len(s.Rows)
}

// BuildGrid converts the scenario rows into a Grid instance.
func (s Scenario) BuildGrid() (Grid, error) {
	if len(s.Rows) == 0 {
		return nil, errors.New("scenario has no rows")
	}
	width := s.Width()
	grid := NewGrid(width, len(s.Rows))
	for y, row := range s.Rows {
		runes := []rune(row)
		if len(runes) != width {
			return nil, fmt.Errorf("row %d has width %d, expected %d", y, len(runes), width)
		}
		for x, r := range runes {
			cell, err := runeToCell(r, s.FireLifetime)
			if err != nil {
				return nil, fmt.Errorf("row %d col %d: %w", y, x, err)
			}
			grid[y][x] = cell
		}
	}
	return grid, nil
}

// GridToScenario creates a Scenario snapshot from the provided grid.
func GridToScenario(name string, grid Grid) Scenario {
	rows := make([]string, len(grid))
	for y, row := range grid {
		runes := make([]rune, len(row))
		for x, cell := range row {
			runes[x] = cellToRune(cell)
		}
		rows[y] = string(runes)
	}
	return Scenario{
		Name: name,
		Rows: rows,
	}
}

// LoadScenario reads a scenario JSON file from disk.
func LoadScenario(path string) (Scenario, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return Scenario{}, err
	}
	var scenario Scenario
	if err := json.Unmarshal(bytes, &scenario); err != nil {
		return Scenario{}, err
	}
	return scenario, nil
}

// SaveScenario writes the scenario to disk using JSON indentation for readability.
func SaveScenario(path string, scenario Scenario) error {
	if len(scenario.Rows) == 0 {
		return errors.New("scenario has no rows to save")
	}
	bytes, err := json.MarshalIndent(scenario, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0o644)
}

func runeToCell(r rune, fireLifetime int) (Cell, error) {
	switch r {
	case '.', '·', 'R':
		return Cell{Type: Rock}, nil
	case 'G':
		return Cell{Type: Grass}, nil
	case 'W':
		return Cell{Type: Water}, nil
	case 'L':
		return Cell{Type: Life}, nil
	case 'F', '🔥':
		cell := NewCell(Fire)
		if fireLifetime > 0 {
			cell.State = fireLifetime
		}
		return cell, nil
	default:
		return Cell{}, fmt.Errorf("unknown cell rune %q", r)
	}
}

func cellToRune(cell Cell) rune {
	switch cell.Type {
	case Rock:
		return '.'
	case Grass:
		return 'G'
	case Water:
		return 'W'
	case Life:
		return 'L'
	case Fire:
		return 'F'
	default:
		return '?'
	}
}
