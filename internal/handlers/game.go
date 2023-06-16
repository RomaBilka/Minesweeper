package handlers

import (
	"encoding/json"
	"github.com/RomaBiliak/Minesweeper/pkg/game"
	"net/http"
)

var g *game.Game
var err error

type cell struct {
	N int `json:"n"`
	M int `json:"m"`
}

type newGameData struct {
	cell
	NumberMines int `json:"numberMines"`
}

func StartGame(w http.ResponseWriter, r *http.Request) {

	gameData := &newGameData{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(gameData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	g, err = game.NewGame(gameData.N, gameData.M, gameData.NumberMines)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(clearDataForResponse(*g)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func OpenCell(w http.ResponseWriter, r *http.Request) {
	c := &cell{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := g.OpenCell(c.N, c.M); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(clearDataForResponse(*g)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DisabledEnabledCell(w http.ResponseWriter, r *http.Request) {
	c := &cell{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := g.DisabledEnabledCell(c.N, c.M); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(clearDataForResponse(*g)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func clearDataForResponse(g game.Game) game.Game {
	cells := make([][]game.Cell, g.N)
	for n := range g.Cells {
		cells[n] = make([]game.Cell, g.M)
		copy(cells[n], g.Cells[n])
	}

	g.Cells = cells
	for n := 0; n < g.N; n++ {
		for m := 0; m < g.M; m++ {
			if !g.Cells[n][m].IsOpen {
				g.Cells[n][m].IsMine = false
				g.Cells[n][m].NumberNeighborhoodMine = 0
			}
		}
	}

	return g
}
