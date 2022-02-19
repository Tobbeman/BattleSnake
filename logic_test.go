package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

const loopNumberForMove = 100

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

	possibleMoves := newPossibleMoves()

	for i := 0; i < loopNumberForMove; i++ {
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
	}

	possibleMoves := newPossibleMoves()

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

		for key, value := range test.pMoves {
			if value != moves[key] {
				t.Errorf("snake does not avoid wall: %+v, %+v", test, moves)
			}
		}
	}
}

func TestSelfAvoidance(t *testing.T) {
	tests := []struct {
		cord   Coord
		body   []Coord
		pMoves map[string]bool
	}{
		{
			cord: Coord{X: 2, Y: 2},
			body: []Coord{{X: 2, Y: 2}, {X: 2, Y: 3}, {X: 1, Y: 3}, {X: 1, Y: 2}, {X: 1, Y: 1}},
			pMoves: map[string]bool{
				"up":    false,
				"down":  true,
				"left":  false,
				"right": true,
			},
		},
	}

	possibleMoves := newPossibleMoves()

	for _, test := range tests {
		me := Battlesnake{
			Head: test.cord,
			Body: test.body,
		}
		state := GameState{
			Board: Board{
				Width:  5,
				Height: 5,
			},
			You: me,
		}

		moves := avoidSelf(state, possibleMoves)

		for key, value := range test.pMoves {
			if value != moves[key] {
				t.Errorf("snake does not avoid self: %+v, %+v", test, moves)
			}
		}
	}
}

func TestSnakeAvoidance(t *testing.T) {
	tests := []struct {
		me     int
		snakes []Battlesnake
		pMoves map[string]bool
	}{
		{
			me: 1,
			snakes: []Battlesnake{
				{Body: []Coord{{X: 1, Y: 2}, {X: 2, Y: 2}, {X: 3, Y: 2}}},
				{Body: []Coord{{X: 1, Y: 1}, {X: 2, Y: 1}}}, // Me
				{Body: []Coord{{X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}}},
			},
			pMoves: map[string]bool{
				"up":    false,
				"down":  false,
				"left":  false,
				"right": true,
			},
		},
	}

	possibleMoves := newPossibleMoves()

	for _, test := range tests {
		me := test.snakes[test.me]
		me.Head = me.Body[1]

		state := GameState{
			Board: Board{
				Width:  5,
				Height: 5,
				Snakes: test.snakes,
			},
			You: me,
		}

		moves, _ := avoidSnakes(state, possibleMoves)
		fail := false
		for key, value := range test.pMoves {
			if value != moves[key] {
				fail = true
			}
		}
		if fail {
			t.Errorf("snake does not avoid self: %+v, %+v", test, moves)
		}
	}
}

func TestMovement(t *testing.T) {
	tests := []struct {
		me     int
		snakes []Battlesnake
		pMoves map[string]bool
	}{
		{
			me: 1,
			snakes: []Battlesnake{
				{Body: []Coord{{X: 1, Y: 2}, {X: 2, Y: 2}, {X: 3, Y: 2}}},
				{Body: []Coord{{X: 1, Y: 1}, {X: 2, Y: 1}}}, // Me
				{Body: []Coord{{X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}}},
			},
			pMoves: map[string]bool{
				"up":    false,
				"down":  false,
				"left":  false,
				"right": true,
			},
		},
	}

	for _, test := range tests {
		me := test.snakes[test.me]
		me.Head = me.Body[1]

		state := GameState{
			Board: Board{
				Width:  5,
				Height: 5,
				Snakes: test.snakes,
			},
			You: me,
		}

		for i := 0; i < loopNumberForMove; i++ {
			move := move(state).Move
			fail := true
			for key, value := range test.pMoves {
				if key == move && value {
					fail = false
				}
			}
			if fail {
				t.Fatalf("got the wrong move: '%s' %v", move, test.pMoves)
			}
		}
	}
}
