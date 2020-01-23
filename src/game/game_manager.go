package game

import (
	"log"
	"time"
)

var zombieMoveDelay = 2 * time.Second

// Run starts a new game and periodically updates the zombie's position.
// Every time a zombie changes position, that position is sent to zombiePositionChannel channel.
// If the zombie reaches the end, the game is finished. Once that happens, gameLostChannel channel will be triggered.
// The game can be stopped by calling close() on the channel that the method returns.
func (game *Game) Run(zombiePositionChannel chan Position, gameLostChannel chan bool) chan struct{} {

	ticker := time.NewTicker(zombieMoveDelay)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				var zombiePosition Position
				if !game.IsGameRunning() {
					zombiePosition = game.StartGame()
				} else {
					zombiePosition = game.MoveZombie()
				}
				zombiePositionChannel <- zombiePosition

				if game.IsGameFinished() {
					log.Println("Game is finished. Stopping.")
					ticker.Stop()
					gameLostChannel <- true
					return
				}
			case <-quit:
				game.StopGame()
				ticker.Stop()
				return
			}
		}
	}()

	return quit
}
