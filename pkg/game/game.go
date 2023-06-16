package game

import (
	"errors"
	"fmt"
	"math/rand"
)

const (
	GameStatusWon  int = 1
	GameStatusLost int = 2
)

// Game
// GameStatus continues=0, won=1, lost=2
type Game struct {
	N            int      `json:"n"`
	M            int      `json:"m"`
	NumbersMines int      `json:"numbersMines"`
	Cells        [][]Cell `json:"cells"`
	GameStatus   int      `json:"gameStatus"`
}

// Cell
// numberNeighborhoodMine - the number of black holes in the neighborhood
type Cell struct {
	NumberNeighborhoodMine int  `json:"numberNeighborhoodMine"`
	IsOpen                 bool `json:"isOpen"`
	IsDisabled             bool `json:"isDisabled"`
	IsMine                 bool `json:"isMine"`
}

func NewGame(n, m, numberMines int) (*Game, error) {
	if n*m < numberMines {
		return nil, errors.New("number of black holes can not be more as n*m")
	}
	if n == 0 || m == 0 || numberMines == 0 {
		return nil, errors.New("n, m or the number of black holes can not be zero")
	}

	g := Game{
		N:            n,
		M:            m,
		NumbersMines: numberMines,
	}
	g.initCells()
	g.setMines()
	g.countNeighborhoodMine()

	return &g, nil
}

func (g *Game) OpenCell(n, m int) error {
	if err := g.validateCellCoordinate(n, m); err != nil {
		return err
	}

	if g.GameStatus != 0 {
		return nil
	}

	if g.Cells[n][m].IsDisabled || g.Cells[n][m].IsOpen {
		return nil
	}

	if g.Cells[n][m].IsMine {
		g.GameStatus = GameStatusLost
		g.openAllMine()
		return nil
	}

	if g.Cells[n][m].NumberNeighborhoodMine > 0 {
		g.Cells[n][m].IsOpen = true
	}

	if g.Cells[n][m].NumberNeighborhoodMine == 0 {
		g.openAllZeroNeighborhood(n, m)
	}

	g.checkIsGameIsWon()
	return nil
}

func (g *Game) DisabledEnabledCell(n, m int) error {
	if err := g.validateCellCoordinate(n, m); err != nil {
		return err
	}

	if g.GameStatus != 0 || g.Cells[n][m].IsOpen {
		return nil
	}

	g.Cells[n][m].IsDisabled = !g.Cells[n][m].IsDisabled

	g.checkIsGameIsWon()
	return nil
}

func (g *Game) validateCellCoordinate(n, m int) error {
	if n < 0 {
		return errors.New("n cannot be less than zero")
	}
	if m < 0 {
		return errors.New("m cannot be less than zero")
	}
	if n > g.N {
		return errors.New(fmt.Sprintf("n cannot be more than %d", g.N))
	}
	if m > g.M {
		return errors.New(fmt.Sprintf("M cannot be more than %d", g.M))
	}

	return nil
}

func (g *Game) initCells() {
	g.Cells = make([][]Cell, g.N, g.N)

	for n := 0; n < g.N; n++ {
		g.Cells[n] = make([]Cell, g.M, g.M)
	}
}

func (g *Game) setMines() {
	blackHole := 0

	for blackHole < g.NumbersMines {
		n := rand.Intn(g.N)
		m := rand.Intn(g.M)
		if !g.Cells[n][m].IsMine {
			g.Cells[n][m].IsMine = true
			blackHole++
		}
	}
}

func (g *Game) countNeighborhoodMine() {
	for n := 0; n < g.N; n++ {
		for m := 0; m < g.M; m++ {
			if g.Cells[n][m].IsMine {
				continue
			}
			g.Cells[n][m].NumberNeighborhoodMine = g.countNeighborhoodMineForOneCell(n, m)
		}
	}
}

func (g *Game) countNeighborhoodMineForOneCell(n, m int) int {
	total := 0

	for i := n - 1; i <= n+1; i++ {
		for j := m - 1; j <= m+1; j++ {
			if i >= 0 && j >= 0 && i < g.N && j < g.M {
				if g.Cells[i][j].IsMine {
					total++
				}
			}
		}
	}

	return total
}

func (g *Game) openAllMine() {
	for n := 0; n < g.N; n++ {
		for m := 0; m < g.M; m++ {
			if g.Cells[n][m].IsMine {
				g.Cells[n][m].IsOpen = true
				g.Cells[n][m].IsDisabled = false
			}
		}
	}
}

func (g *Game) openAllZeroNeighborhood(n, m int) {
	if !g.Cells[n][m].IsDisabled && !g.Cells[n][m].IsOpen && !g.Cells[n][m].IsMine {
		g.Cells[n][m].IsOpen = true
	} else {
		return
	}

	if g.Cells[n][m].NumberNeighborhoodMine > 0 {
		return
	}

	for i := n - 1; i <= n+1; i++ {
		for j := m - 1; j <= m+1; j++ {
			if i >= 0 && j >= 0 && i < g.N && j < g.M {
				g.openAllZeroNeighborhood(i, j)
			}
		}
	}
}

func (g *Game) checkIsGameIsWon() {
	totalClosed := 0
	totalDisabled := 0

	for n := 0; n < g.N; n++ {
		for m := 0; m < g.M; m++ {
			if !g.Cells[n][m].IsOpen != g.Cells[n][m].IsMine {
				totalClosed++
				return
			}
			if g.Cells[n][m].IsDisabled {
				totalDisabled++
				return
			}
		}
	}

	if totalClosed == 0 && totalDisabled == 0 {
		g.GameStatus = GameStatusWon
	}
}
