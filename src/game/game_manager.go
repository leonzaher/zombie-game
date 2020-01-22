package game

import (
	"log"
	"time"
)

const zombieMoveDelay = 2 * time.Second

func (game *Game) Run(sender chan Position) chan struct{} {

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
				sender <- zombiePosition

				if game.IsGameFinished() {
					log.Println("Game is finished. Stopping.")
					ticker.Stop()
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
