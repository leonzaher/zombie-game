package game

import (
	"testing"
	"time"
)

func TestGame_Run_UntilEnd(t *testing.T) {
	gameInstance := Game{}
	zombiePositionChannel := make(chan Position)
	gameLostChannel := make(chan bool)

	zombieMoveDelay = 10 * time.Millisecond

	gameInstance.Run(zombiePositionChannel, gameLostChannel)

	for {
		select {
		case position := <-zombiePositionChannel:
			println(position.ToString())
		case gameFinished := <-gameLostChannel:
			if !gameFinished {
				t.Errorf("gameFinished should contain true")
			}
			println("Game lost, all good")
			return
		case <-time.After(zombieMoveDelay * (maxColumnPosition + 1)):
			t.Errorf("Game should have been lost by now, but it isn't")
		}
	}
}

func TestGame_Run_GameStopped(t *testing.T) {
	gameInstance := Game{}
	zombiePositionChannel := make(chan Position)
	gameLostChannel := make(chan bool)

	quitGame := gameInstance.Run(zombiePositionChannel, gameLostChannel)

	time.Sleep(10 * time.Millisecond)

	close(quitGame)
}
