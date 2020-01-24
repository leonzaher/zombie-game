package game

import (
	"fmt"
	"log"
	"time"
	"zombie-game/src/chTypes"
)

var zombieMoveDelay = 2 * time.Second

// Run starts a new game and periodically updates the zombie's position.
// Every time a zombie changes position, that position is sent to zombiePositionChannel channel.
// If the zombie reaches the end, the game is finished. Once that happens, gameLostChannel channel will be triggered.
// The game can be stopped by calling close() on the channel that the method returns.
func (game *Game) Run(broadcast chan chTypes.Broadcast, shotAttempt chan chTypes.ShotAttempt, gameStop chan string) {

	ticker := time.NewTicker(zombieMoveDelay)
	exit := make(chan struct{})
	var exitMessage string
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

				zombiePosMessage := fmt.Sprintf("WALK %s %d %d",
					game.GetZombieName(), zombiePosition.Row, zombiePosition.Column)
				broadcast <- chTypes.Broadcast{Message: zombiePosMessage, HostingPlayerName: game.HostingPlayerName}

				if game.IsGameFinished() {
					log.Println("Game is finished. Stopping.")
					exitMessage = fmt.Sprintf("Game lost, %s reached the end!", game.GetZombieName())
					close(exit)
				}

			case shotAttempt := <-shotAttempt:
				if game.IsZombieOnPosition(shotAttempt.Row, shotAttempt.Column) {
					exitMessage = fmt.Sprintf("BOOM %s wins!", shotAttempt.PlayerName)
					close(exit)
				} else {
					missMessage := fmt.Sprintf("MISS %s at [%d %d]",
						shotAttempt.PlayerName, shotAttempt.Row, shotAttempt.Column)
					broadcast <- chTypes.Broadcast{Message: missMessage, HostingPlayerName: game.HostingPlayerName}
				}

			case <-exit:
				broadcast <- chTypes.Broadcast{Message: exitMessage, HostingPlayerName: game.HostingPlayerName}
				gameStop <- game.HostingPlayerName
				ticker.Stop()
				game.StopGame()
				return
			}
		}
	}()
}
