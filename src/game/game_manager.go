package game

import (
	"log"
	"time"
)

const zombieMoveDelay = 2 * time.Second

func Run(game *Game) chan struct{} {

	ticker := time.NewTicker(zombieMoveDelay)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				MoveZombie(game)
				if IsGameFinished(*game) {
					log.Println("Game is finished. Stopping.")
					ticker.Stop()
					return
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return quit
}
