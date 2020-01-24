package game

import (
	"testing"
	"time"
	"zombie-game/src/chTypes"
)

func TestGame_Run_UntilEnd(t *testing.T) {
	playerName := "test"
	gameInstance := Game{HostingPlayerName: playerName}

	broadcast := make(chan chTypes.Broadcast)
	shotAttempt := make(chan chTypes.ShotAttempt)
	stopGame := make(chan string)

	zombieMoveDelay = 10 * time.Millisecond

	gameInstance.Run(broadcast, shotAttempt, stopGame)

	for {
		select {
		case message := <-broadcast:
			println(message.Message)
		case player := <-stopGame:
			if player != playerName {
				t.Errorf("player should be %s, but is %s", playerName, player)
			}
			println("Game lost, all good")
			return
		case <-time.After(zombieMoveDelay * (maxColumnPosition + 1)):
			t.Errorf("Game should have been lost by now, but it isn't")
		}
	}
}
