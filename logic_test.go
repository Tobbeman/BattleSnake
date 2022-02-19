package main

import (
	"testing"
)

func TestNeckAvoidance(t *testing.T) {
	// Arrange
	me := Battlesnake{
		// Length 3, facing right
		Head: Coord{X: 2, Y: 0},
		Body: []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
	}
	state := GameState{
		Board: Board{
			Snakes: []Battlesnake{me},
		},
		You: me,
	}

	possibleMoves := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	// Act 1,000x (this isn't a great way to test, but it's okay for starting out)
	for i := 0; i < 1000; i++ {
		moves := avoidNeck(state, possibleMoves)
		// Assert never move left
		if moves["left"] == true {
			t.Errorf("snake moved onto its own neck")
		}
	}
}

func TestWallAvoidance(t *testing.T) {
	tests := []struct {
		cord   Coord
		pMoves map[string]bool
	}{
		{
			cord: Coord{X: 0, Y: 0},
			pMoves: map[string]bool{
				"up":    true,
				"down":  false,
				"left":  false,
				"right": true,
			},
		},
		{
			cord: Coord{X: 0, Y: 4},
			pMoves: map[string]bool{
				"up":    false,
				"down":  true,
				"left":  false,
				"right": true,
			},
		},
		{
			cord: Coord{X: 4, Y: 4},
			pMoves: map[string]bool{
				"up":    false,
				"down":  true,
				"left":  true,
				"right": false,
			},
		},
		{
			cord: Coord{X: 4, Y: 0},
			pMoves: map[string]bool{
				"up":    true,
				"down":  false,
				"left":  true,
				"right": false,
			},
		},
		{
			cord: Coord{X: 3, Y: 0},
			pMoves: map[string]bool{
				"up":    true,
				"down":  false,
				"left":  true,
				"right": true,
			},
		},
	}

	possibleMoves := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	for _, test := range tests {
		me := Battlesnake{
			Head: test.cord,
		}
		state := GameState{
			Board: Board{
				Width:  5,
				Height: 5,
			},
			You: me,
		}

		moves := avoidWall(state, possibleMoves)

		for key, value := range possibleMoves {
			if value != moves[key] {
				t.Errorf("snake does not avoid wall: %+v", test)
			}
		}
	}
}
