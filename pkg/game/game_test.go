package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	testCases := []struct {
		name             string
		n                int
		m                int
		numberBlackHoles int
		error            string
	}{
		{
			name:             "numberBlackHoles > n*n",
			n:                1,
			m:                1,
			numberBlackHoles: 2,
			error:            "number of black holes can not be more as n*m",
		},
		{
			name:             "numberBlackHoles is zero",
			n:                1,
			m:                1,
			numberBlackHoles: 0,
			error:            "n, m or the number of black holes can not be zero",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := NewGame(testCase.n, testCase.m, testCase.numberBlackHoles)
			assert.Errorf(t, err, testCase.error)
		})
	}
}

func TestValidateCellCoordinate(t *testing.T) {
	g, _ := NewGame(8, 8, 10)
	testCases := []struct {
		name  string
		n     int
		m     int
		error string
	}{
		{name: "n cannot be less than zero", n: -1, m: 1, error: "n cannot be less than zero"},
		{name: "m cannot be less than zero", n: 1, m: -1, error: "m cannot be less than zero"},
		{name: "n cannot be more than 8", n: 9, m: 1, error: "n cannot be more than 8"},
		{name: "m cannot be more than 8", n: 1, m: 9, error: "m cannot be more than 8"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := g.validateCellCoordinate(testCase.n, testCase.m)
			assert.Errorf(t, err, testCase.error)
		})
	}
}

func TestCheckIsGameIsWon(t *testing.T) {
	testCases := []struct {
		name       string
		g          *Game
		gameStatus int
	}{
		{
			name: "Game continues",
			g: &Game{
				Cells: [][]Cell{
					{{IsOpen: true}, {IsBlackHole: true}},
					{{IsOpen: true}, {}},
				},
				N:                 2,
				M:                 2,
				NumbersBlackHoles: 1,
			},
			gameStatus: 0,
		},
		{
			name: "Game won",
			g: &Game{
				Cells: [][]Cell{
					{{IsOpen: true}, {IsBlackHole: true}},
					{{IsOpen: true}, {IsOpen: true}},
				},
				N:                 2,
				M:                 2,
				NumbersBlackHoles: 1,
			},
			gameStatus: 1,
		},
		{
			name: "Game with disabled cell",
			g: &Game{
				Cells: [][]Cell{
					{{IsOpen: true}, {IsBlackHole: true, IsDisabled: true}},
					{{IsOpen: true}, {IsOpen: true}},
				},
				N:                 2,
				M:                 2,
				NumbersBlackHoles: 0,
			},
			gameStatus: 0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.g.checkIsGameIsWon()
			assert.Equal(t, testCase.gameStatus, testCase.g.GameStatus)
		})
	}
}

func TestOpenAllBlackHole(t *testing.T) {
	g := &Game{
		Cells: [][]Cell{
			{{}, {IsBlackHole: true}},
			{{IsBlackHole: true}, {}},
		},
		N:                 2,
		M:                 2,
		NumbersBlackHoles: 1,
	}

	testCases := []struct {
		name   string
		n, m   int
		isOpen bool
	}{
		{name: "0 0 close", n: 0, m: 0, isOpen: false},
		{name: "0 1 open", n: 0, m: 1, isOpen: true},
		{name: "1 0 open", n: 1, m: 0, isOpen: true},
		{name: "1 1 close", n: 1, m: 1, isOpen: false},
	}

	g.openAllBlackHole()
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.isOpen, g.Cells[testCase.n][testCase.m].IsOpen)
		})
	}
}

func TestDisabledEnabledCell(t *testing.T) {
	g := &Game{
		Cells: [][]Cell{
			{{IsDisabled: true}, {}},
			{{IsOpen: true}, {}},
		},
		N: 2,
		M: 2,
	}

	testCases := []struct {
		name       string
		n, m       int
		isDisabled bool
	}{
		{name: "0 0 Enabled", n: 0, m: 0, isDisabled: false},
		{name: "1 1 Disabled", n: 1, m: 1, isDisabled: true},
		{name: "1 1 Enabled", n: 1, m: 0, isDisabled: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := g.DisabledEnabledCell(testCase.n, testCase.m)
			assert.NoError(t, err)
			assert.Equal(t, testCase.isDisabled, g.Cells[testCase.n][testCase.m].IsDisabled)
		})
	}
}

func TestOpenAllZeroNeighborhood(t *testing.T) {
	g := &Game{
		Cells: [][]Cell{
			{{}, {}, {}},
			{{}, {NumberNeighborhoodBlackHole: 1}, {NumberNeighborhoodBlackHole: 1}},
			{{}, {NumberNeighborhoodBlackHole: 1}, {IsBlackHole: true}},
		},
		N:                 3,
		M:                 3,
		NumbersBlackHoles: 1,
	}
	testCases := []struct {
		name   string
		n, m   int
		isOpen bool
	}{
		{name: "0 0 Open", n: 0, m: 0, isOpen: true},
		{name: "0 1 Open", n: 0, m: 1, isOpen: true},
		{name: "0 2 Open", n: 0, m: 2, isOpen: true},
		{name: "1 0 Open", n: 1, m: 0, isOpen: true},
		{name: "1 1 Open", n: 1, m: 1, isOpen: true},
		{name: "1 2 Open", n: 1, m: 2, isOpen: true},
		{name: "2 0 Open", n: 2, m: 0, isOpen: true},
		{name: "2 1 Open", n: 2, m: 1, isOpen: true},
		{name: "2 2 Open", n: 2, m: 2, isOpen: false},
	}
	g.openAllZeroNeighborhood(0, 0)
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.isOpen, g.Cells[testCase.n][testCase.m].IsOpen)
		})
	}
}

func TestOpenCell(t *testing.T) {
	testCases := []struct {
		name                        string
		n, m                        int
		isOpen                      bool
		numberNeighborhoodBlackHole int
		gameStatus                  int
	}{
		{name: "0 0 Open", n: 0, m: 0, isOpen: true, gameStatus: 1},
		{name: "1 1 Open", n: 1, m: 1, isOpen: true, numberNeighborhoodBlackHole: 1},
		{name: "2 2 Open", n: 2, m: 2, isOpen: true, gameStatus: 2},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			g := &Game{
				Cells: [][]Cell{
					{{}, {}, {}},
					{{}, {NumberNeighborhoodBlackHole: 1}, {NumberNeighborhoodBlackHole: 1}},
					{{}, {NumberNeighborhoodBlackHole: 1}, {IsBlackHole: true}},
				},
				N:                 3,
				M:                 3,
				NumbersBlackHoles: 1,
			}
			err := g.OpenCell(testCase.n, testCase.m)
			if err != nil {

			}
			assert.NoError(t, err)
			assert.Equal(t, testCase.isOpen, g.Cells[testCase.n][testCase.m].IsOpen)
			assert.Equal(t, testCase.numberNeighborhoodBlackHole, g.Cells[testCase.n][testCase.m].NumberNeighborhoodBlackHole)
			assert.Equal(t, testCase.gameStatus, g.GameStatus)
		})
	}
}
