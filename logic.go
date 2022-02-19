package main

// This file can be a nice home for your Battlesnake logic and related helper functions.
//
// We have started this for you, with a function to help remove the 'neck' direction
// from the list of possible moves!

import (
	"log"
	"math"
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
	possibleMoves := newPossibleMoves()
	preferedMoves := newPossibleMoves()

	// Step 0: Don't let your Battlesnake move back in on it's own neck
	// possibleMoves = avoidNeck(state, possibleMoves)

	// Step 1 - Don't hit walls.
	possibleMoves = avoidWall(state, possibleMoves)

	// Step 2 - Don't hit yourself.
	// possibleMoves = avoidSelf(state, possibleMoves)

	// Step 3 - Don't collide with others. Or yourself!
	possibleMoves, preferedMoves = avoidSnakes(state, possibleMoves)

	// TODO: Step 4 - Find food.
	// Use information in GameState to seek out and find food.

	// Finally, choose a move from the available safe moves.
	// TODO: Step 5 - Select a move to make based on strategy, rather than random.
	var nextMove string
	safeMoves := []string{}
	moveMap := possibleMoves // Reference!

	for _, isSafe := range preferedMoves {
		if isSafe {
			moveMap = preferedMoves
		}
	}

	for move, isSafe := range moveMap {
		if isSafe {
			safeMoves = append(safeMoves, move)
		}
	}

	if len(safeMoves) == 0 {
		nextMove = "down"
		log.Printf("%s MOVE %d: No safe moves detected! Moving %s\n", state.Game.ID, state.Turn, nextMove)
	} else {
		nextMove = getFood(state, safeMoves)
		if nextMove == "" {
			nextMove = safeMoves[rand.Intn(len(safeMoves))]
		}
		log.Printf("%s MOVE %d: %s\n", state.Game.ID, state.Turn, nextMove)
	}
	return BattlesnakeMoveResponse{
		Move: nextMove,
	}
}

func newPossibleMoves() map[string]bool {
	return map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
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

func avoidSnakes(state GameState, possibleMoves map[string]bool) (map[string]bool, map[string]bool) {
	_possibleMoves := copyMoves(possibleMoves)
	_preferedMoves := copyMoves(possibleMoves)

	myHead := state.You.Head

	for _, snake := range state.Board.Snakes {
		for _, link := range snake.Body {
			if myHead.Y == link.Y {
				if myHead.X+1 == link.X {
					_possibleMoves["right"] = false
					_preferedMoves["right"] = false
				}
				if myHead.X-1 == link.X {
					_possibleMoves["left"] = false
					_preferedMoves["left"] = false
				}
				if myHead.X+2 == link.X {
					_preferedMoves["right"] = false
				}
				if myHead.X-2 == link.X {
					_preferedMoves["left"] = false
				}
			}
			if myHead.X == link.X {
				if myHead.Y+1 == link.Y {
					_possibleMoves["up"] = false
					_preferedMoves["up"] = false
				}
				if myHead.Y-1 == link.Y {
					_possibleMoves["down"] = false
					_preferedMoves["down"] = false
				}
				if myHead.Y+2 == link.Y {
					_preferedMoves["up"] = false
				}
				if myHead.Y-2 == link.Y {
					_preferedMoves["down"] = false
				}
			}
		}
	}

	return _possibleMoves, _preferedMoves
}

func getFood(state GameState, safeMoves []string) string {
	move := ""

	// Lets create a ray!
	myHead := state.You.Head

	var closet Coord
	closetDist := state.Board.Height * state.Board.Width

	for _, food := range state.Board.Food {
		if getDistance(myHead, food) < closetDist {
			closet = food
		}
	}

	for _, dir := range safeMoves {
		switch dir {
		case "up":
			if (getDistance(Coord{myHead.X, myHead.Y + 1}, closet) > getDistance(myHead, closet)) {
				move = "up"
			}
		case "down":
			if (getDistance(Coord{myHead.X, myHead.Y - 1}, closet) > getDistance(myHead, closet)) {
				move = "down"
			}
		case "left":
			if (getDistance(Coord{myHead.X - 1, myHead.Y}, closet) > getDistance(myHead, closet)) {
				move = "left"
			}
		case "right":
			if (getDistance(Coord{myHead.X + 1, myHead.Y}, closet) > getDistance(myHead, closet)) {
				move = "right"
			}
		}

	}

	return move
}

func getDistance(start, end Coord) int {
	return int(math.Abs(float64(start.Y)-float64(end.Y)) + math.Abs(float64(start.X)-float64(end.X)))
}
