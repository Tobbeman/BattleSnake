package main

// This file can be a nice home for your Battlesnake logic and related helper functions.
//
// We have started this for you, with a function to help remove the 'neck' direction
// from the list of possible moves!

import (
	"log"
	"math/rand"
)

// This function is called when you register your Battlesnake on play.battlesnake.com
// See https://docs.battlesnake.com/guides/getting-started#step-4-register-your-battlesnake
// It controls your Battlesnake appearance and author permissions.
// For customization options, see https://docs.battlesnake.com/references/personalization
// TIP: If you open your Battlesnake URL in browser you should see this data.
func info() BattlesnakeInfoResponse {
	log.Println("INFO")
	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "Tobbeman",
		Color:      "#4C89C8",
		Head:       "shades",
		Tail:       "default", // TODO: Personalize
	}
}

// This function is called everytime your Battlesnake is entered into a game.
// The provided GameState contains information about the game that's about to be played.
// It's purely for informational purposes, you don't have to make any decisions here.
func start(state GameState) {
	log.Printf("%s START\n", state.Game.ID)
}

// This function is called when a game your Battlesnake was in has ended.
// It's purely for informational purposes, you don't have to make any decisions here.
func end(state GameState) {
	log.Printf("%s END\n\n", state.Game.ID)
}

// This function is called on every turn of a game. Use the provided GameState to decide
// where to move -- valid moves are "up", "down", "left", or "right".
// We've provided some code and comments to get you started.
func move(state GameState) BattlesnakeMoveResponse {
	possibleMoves := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	log.Printf("%+v", state.You.Body)

	// Step 0: Don't let your Battlesnake move back in on it's own neck
	// possibleMoves = avoidNeck(state, possibleMoves)

	// Step 1 - Don't hit walls.
	possibleMoves = avoidWall(state, possibleMoves)

	// Step 2 - Don't hit yourself.
	// possibleMoves = avoidSelf(state, possibleMoves)

	// Step 3 - Don't collide with others.
	possibleMoves = avoidSnakes(state, possibleMoves)

	// TODO: Step 4 - Find food.
	// Use information in GameState to seek out and find food.

	// Finally, choose a move from the available safe moves.
	// TODO: Step 5 - Select a move to make based on strategy, rather than random.
	var nextMove string

	safeMoves := []string{}
	for move, isSafe := range possibleMoves {
		if isSafe {
			safeMoves = append(safeMoves, move)
		}
	}

	if len(safeMoves) == 0 {
		nextMove = "down"
		log.Printf("%s MOVE %d: No safe moves detected! Moving %s\n", state.Game.ID, state.Turn, nextMove)
	} else {
		nextMove = safeMoves[rand.Intn(len(safeMoves))]
		log.Printf("%s MOVE %d: %s\n", state.Game.ID, state.Turn, nextMove)
	}
	return BattlesnakeMoveResponse{
		Move: nextMove,
	}
}

func copyMoves(input map[string]bool) map[string]bool {
	ret := map[string]bool{}
	for k, v := range input {
		ret[k] = v
	}
	return ret
}

func avoidNeck(state GameState, possibleMoves map[string]bool) map[string]bool {
	moves := copyMoves(possibleMoves)

	myHead := state.You.Head
	myNeck := state.You.Body[1]
	if myNeck.X < myHead.X {
		moves["left"] = false
	} else if myNeck.X > myHead.X {
		moves["right"] = false
	} else if myNeck.Y < myHead.Y {
		moves["down"] = false
	} else if myNeck.Y > myHead.Y {
		moves["up"] = false
	}

	return moves
}

func avoidWall(state GameState, possibleMoves map[string]bool) map[string]bool {
	moves := copyMoves(possibleMoves)

	myHead := state.You.Head
	boardWidth := state.Board.Width
	boardHeight := state.Board.Height
	if myHead.X+1 > boardWidth-1 {
		moves["right"] = false
	}
	if myHead.X-1 < 0 {
		moves["left"] = false
	}
	if myHead.Y+1 > boardHeight-1 {
		moves["up"] = false
	}
	if myHead.Y-1 < 0 {
		moves["down"] = false
	}

	return moves
}

func avoidSelf(state GameState, possibleMoves map[string]bool) map[string]bool {
	moves := copyMoves(possibleMoves)

	myHead := state.You.Head
	myBody := state.You.Body

	for _, link := range myBody {
		if myHead.X+1 == link.X && myHead.Y == link.Y {
			moves["right"] = false
		}
		if myHead.X-1 == link.X && myHead.Y == link.Y {
			moves["left"] = false
		}
		if myHead.Y+1 == link.Y && myHead.X == link.X {
			moves["up"] = false
		}
		if myHead.Y-1 == link.Y && myHead.X == link.X {
			moves["down"] = false
		}
	}

	return moves
}

func avoidSnakes(state GameState, possibleMoves map[string]bool) map[string]bool {
	moves := copyMoves(possibleMoves)

	myHead := state.You.Head

	for _, snake := range state.Board.Snakes {
		for _, link := range snake.Body {
			if myHead.Y == link.Y {
				if myHead.X+1 == link.X {
					moves["right"] = false
				}
				if myHead.X-1 == link.X {
					moves["left"] = false
				}
			}
			if myHead.X == link.X {
				if myHead.Y+1 == link.Y {
					moves["up"] = false
				}
				if myHead.Y-1 == link.Y {
					moves["down"] = false
				}
			}
		}
	}

	return moves
}
