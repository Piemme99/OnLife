package main

import (
	"OnLife/world"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	scenarioPath := flag.String("scenario", "", "Path to a scenario JSON file")
	flag.Parse()

	grid, err := loadInitialGrid(*scenarioPath)
	if err != nil {
		log.Fatalf("load grid: %v", err)
	}

	sim := NewSimulation(grid)
	hub := NewHub()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(sim, hub, w, r)
	})

	log.Printf("WebSocket server listening on %s (scenario: %s)", *addr, scenarioLabel(*scenarioPath))
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}

func scenarioLabel(path string) string {
	if path == "" {
		return "random"
	}
	return path
}

type Simulation struct {
	mu   sync.Mutex
	grid *world.Grid
	tick int
}

func NewSimulation(grid *world.Grid) *Simulation {
	return &Simulation{grid: grid}
}

func (s *Simulation) Tick() StateMessage {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.grid.Update()
	s.tick++
	return s.stateMessageLocked()
}

func (s *Simulation) Snapshot() StateMessage {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stateMessageLocked()
}

func (s *Simulation) stateMessageLocked() StateMessage {
	return StateMessage{
		Type: "state",
		Tick: s.tick,
		Grid: gridToSymbols(*s.grid),
	}
}

type StateMessage struct {
	Type string     `json:"type"`
	Tick int        `json:"tick"`
	Grid [][]string `json:"grid"`
}

type ErrorMessage struct {
	Type  string `json:"type"`
	Error string `json:"error"`
}

type Command struct {
	Type string `json:"type"`
}

type Hub struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]struct{}
}

func NewHub() *Hub {
	return &Hub{clients: make(map[*websocket.Conn]struct{})}
}

func (h *Hub) Add(conn *websocket.Conn) {
	h.mu.Lock()
	h.clients[conn] = struct{}{}
	h.mu.Unlock()
}

func (h *Hub) Remove(conn *websocket.Conn) {
	h.mu.Lock()
	delete(h.clients, conn)
	h.mu.Unlock()
	conn.Close()
}

func (h *Hub) Broadcast(msg interface{}) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("broadcast marshal: %v", err)
		return
	}
	h.mu.Lock()
	for conn := range h.clients {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("broadcast write: %v", err)
			conn.Close()
			delete(h.clients, conn)
		}
	}
	h.mu.Unlock()
}

func serveWS(sim *Simulation, hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade: %v", err)
		return
	}
	hub.Add(conn)

	if err := conn.WriteJSON(sim.Snapshot()); err != nil {
		log.Printf("send snapshot: %v", err)
		hub.Remove(conn)
		return
	}

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			hub.Remove(conn)
			return
		}
		var cmd Command
		if err := json.Unmarshal(data, &cmd); err != nil {
			writeError(conn, fmt.Sprintf("invalid command: %v", err))
			continue
		}
		switch cmd.Type {
		case "tick":
			state := sim.Tick()
			hub.Broadcast(state)
		case "state":
			if err := conn.WriteJSON(sim.Snapshot()); err != nil {
				log.Printf("send snapshot: %v", err)
			}
		default:
			writeError(conn, "unknown command")
		}
	}
}

func writeError(conn *websocket.Conn, msg string) {
	if err := conn.WriteJSON(ErrorMessage{Type: "error", Error: msg}); err != nil {
		log.Printf("send error: %v", err)
	}
}

func loadInitialGrid(path string) (*world.Grid, error) {
	if path == "" {
		w := world.NewWord(world.GridSize, world.GridSize)
		grid := w.Grid
		grid.InitializeRandomGrid()
		return &grid, nil
	}
	scenario, err := world.LoadScenario(path)
	if err != nil {
		return nil, err
	}
	grid, err := scenario.BuildGrid()
	if err != nil {
		return nil, err
	}
	return &grid, nil
}

func gridToSymbols(grid world.Grid) [][]string {
	rows := make([][]string, len(grid))
	for y, row := range grid {
		rows[y] = make([]string, len(row))
		for x, cell := range row {
			rows[y][x] = cellToSymbol(cell)
		}
	}
	return rows
}

func cellToSymbol(cell world.Cell) string {
	switch cell.Type {
	case world.Rock:
		return "."
	case world.Grass:
		return "G"
	case world.Water:
		return "W"
	case world.Life:
		return "L"
	case world.Fire:
		return "F"
	default:
		return "?"
	}
}
