package game

import (
	"log"
	"math/rand"
	"time"
)

const maxRowPosition = 10 - 1
const maxColumnPosition = 30 - 1

type Game struct {
	playerName string
	zombieName string

	zombiePosition position
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewGame(playerName string, zombieName string) Game {
	zombiePosition := generateStartingPosition()
	log.Println("Started new game.", zombieName, "position is ", toString(zombiePosition))
	return Game{
		playerName:     playerName,
		zombieName:     zombieName,
		zombiePosition: zombiePosition,
	}
}

// Zombie always starts along the left side of the board
// row = random, column = 0
func generateStartingPosition() position {
	rows := randRange(0, maxRowPosition)
	return newPosition(rows, 0)
}

func IsZombieOnPosition(game Game, row int, column int) bool {
	return game.zombiePosition.row == row && game.zombiePosition.column == column
}

func MoveZombie(game *Game) {
	game.zombiePosition.row += randRange(-1, 1)
	game.zombiePosition.column += 1

	if game.zombiePosition.row > maxRowPosition {
		game.zombiePosition.row = maxRowPosition
	}
	if game.zombiePosition.row < 0 {
		game.zombiePosition.row = 0
	}

	log.Println(game.zombieName, "is moving to", toString(game.zombiePosition))
}

func IsGameFinished(game Game) bool {
	return game.zombiePosition.column == maxColumnPosition
}

// returns random number in inclusive interval: [min, max]
func randRange(min int, max int) int {
	return rand.Intn(max-min+1) + min
}
