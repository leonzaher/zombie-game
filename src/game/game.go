package game

import (
	"log"
	"math/rand"
	"time"
)

const maxRowPosition = 10 - 1
const maxColumnPosition = 30 - 1

type Game struct {
	HostingPlayerName string
	zombieName        string
	zombiePosition    Position
	gameRunning       bool
}

var zombieNames = [...]string{"Night king", "Frankenstein", "Alien"}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (game *Game) StartGame() Position {
	game.gameRunning = true
	game.zombieName = zombieNames[randRange(0, len(zombieNames)-1)]
	game.zombiePosition = generateStartingPosition()
	log.Println("Started new game.", game.zombieName, "'s position is ", game.zombiePosition.ToString())
	return game.zombiePosition
}

// Zombie always starts along the left side of the board
// row = random, column = 0
func generateStartingPosition() Position {
	rows := randRange(0, maxRowPosition)
	return newPosition(rows, 0)
}

func (game Game) IsGameRunning() bool {
	return game.gameRunning
}

func (game Game) IsZombieOnPosition(row int, column int) bool {
	return game.zombiePosition.Row == row && game.zombiePosition.Column == column
}

func (game Game) GetZombieName() string {
	return game.zombieName
}

func (game *Game) MoveZombie() Position {
	game.zombiePosition.Row += randRange(-1, 1)
	game.zombiePosition.Column += 1

	if game.zombiePosition.Row > maxRowPosition {
		game.zombiePosition.Row = maxRowPosition
	}
	if game.zombiePosition.Row < 0 {
		game.zombiePosition.Row = 0
	}

	log.Println(game.zombieName, "is moving to", game.zombiePosition.ToString())
	return game.zombiePosition
}

func (game Game) IsGameFinished() bool {
	return game.zombiePosition.Column == maxColumnPosition
}

func (game *Game) StopGame() {
	log.Println("Stopping game")
	*game = Game{}
}

// returns random number in inclusive interval: [min, max]
func randRange(min int, max int) int {
	return rand.Intn(max-min+1) + min
}
