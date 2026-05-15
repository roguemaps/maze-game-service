package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var players = make(map[string]Player)
var mazes = make(map[string]Maze)

type Player struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Position  Position `json:"position"`
	Direction string   `json:"direction"`
	Items     []string `json:"items"`
}

type item struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Position struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type Direction string

const (
	North Direction = "north"
	East  Direction = "east"
	South Direction = "south"
	West  Direction = "west"
)

type Wall struct {
	Direction Direction `json:"direction"`
}

type Door struct {
	Direction Direction `json:"direction"`
}

type Cell struct {
	Id          string   `json:"id"`
	Description string   `json:"description"`
	Items       []string `json:"items"`
	Walls       []Wall   `json:"walls"`
	Visited     bool     `json:"visited"`
	Exits       []string `json:"exits"`
	Finished    bool     `json:"finished"`
}

type Maze struct {
	Id      string    `json:"id"`
	Visited []string  `json:"visited"`
	Current Position  `json:"current"`
	Player  Player    `json:"player"`
	OwnerId string    `json:"playerId"`
	Grid    [][]*Cell `json:"grid"`
}

func newMaze(player Player) Maze {
	startingCell := Cell{
		Id:          "cell-0-0",
		Description: "Starting cell",
		Items:       []string{},
	}

	return Maze{
		Id:      fmt.Sprintf("maze-%d", time.Now().UnixNano()),
		Visited: []string{startingCell.Id},
		Current: Position{Row: 0, Col: 0},
		Player:  player,
		OwnerId: player.Id,
		Grid:    [][]*Cell{{&startingCell}},
	}
}

func newPlayer(name string, email string) Player {
	return Player{
		Id:        fmt.Sprintf("player-%d", time.Now().UnixNano()),
		Name:      name,
		Email:     email,
		Position:  Position{Row: 0, Col: 0},
		Direction: "north",
		Items:     []string{},
	}
}

func newMazeHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		PlayerId string `json:"playerId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "unable to decode request body", http.StatusBadRequest)
		return
	}

	if body.PlayerId == "" {
		http.Error(w, "playerId is required", http.StatusBadRequest)
		return
	}

	player, ok := players[body.PlayerId]
	if !ok {
		http.Error(w, "player not found", http.StatusNotFound)
		return
	}

	maze := newMaze(player)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(maze)
}

func newPlayerHandler(w http.ResponseWriter, r *http.Request) {
	var player Player

	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		http.Error(w, "unable to decode player Object", http.StatusBadRequest)
		return
	}

	if player.Name == "" {
		http.Error(w, "player name is required", http.StatusBadRequest)
		return
	}

	if player.Email == "" {
		http.Error(w, "player email is required", http.StatusBadRequest)
		return
	}

	player = newPlayer(player.Name, player.Email)
	players[player.Id] = player

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

func saveMazeHandler(w http.ResponseWriter, r *http.Request) {
	var maze Maze

	if err := json.NewDecoder(r.Body).Decode(&maze); err != nil {
		http.Error(w, "unable to decode maze", http.StatusBadRequest)
		return
	}

	mazes[maze.Id] = maze

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(maze)
}

func getMazeHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	maze, ok := mazes[id]
	if !ok {
		http.Error(w, "maze not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(maze)
}

func getMazesHandler(w http.ResponseWriter, r *http.Request) {

	list := make([]Maze, 0, len(mazes))
	for _, m := range mazes {
		list = append(list, m)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func getPlayerHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	player, ok := players[id]
	if !ok {
		http.Error(w, "player not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

func getPlayersHandler(w http.ResponseWriter, r *http.Request) {
	list := make([]Player, 0, len(players))
	for _, p := range players {
		list = append(list, p)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func main() {
	http.HandleFunc("/newmaze", newMazeHandler)
	http.HandleFunc("/savemaze", saveMazeHandler)
	http.HandleFunc("/getmaze/:id", getMazeHandler)
	http.HandleFunc("/mazes", getMazesHandler)
	http.HandleFunc("/newmazeplayer", newPlayerHandler)
	http.HandleFunc("/getplayer/:id", getPlayerHandler)
	http.HandleFunc("/players", getPlayersHandler)

	log.Println("maze service listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
