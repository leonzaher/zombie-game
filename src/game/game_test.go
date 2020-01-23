package game

import (
	"testing"
)

func TestGame_StartGame(t *testing.T) {
	gameInstance := Game{}
	gameInstance.StartGame()

	if !gameInstance.IsGameRunning() {
		t.Errorf("Game should be running")
	}
	if gameInstance.GetZombieName() == "" {
		t.Errorf("Zombie name cannot be empty")
	}
}

func TestGame_IsGameRunning(t *testing.T) {
	gameInstance := Game{}
	if gameInstance.IsGameRunning() {
		t.Errorf("Game should NOT be running before start")
	}

	gameInstance.StartGame()
	if !gameInstance.IsGameRunning() {
		t.Errorf("Game should be running after start")
	}

	gameInstance.StopGame()
	if gameInstance.IsGameRunning() {
		t.Errorf("Game should NOT be running after stop")
	}
}

func TestGame_IsZombieOnPosition(t *testing.T) {
	gameInstance := Game{}
	gameInstance.StartGame()

	zombieRow := -100
	zombieCol := 100
	gameInstance.zombiePosition = Position{Row: zombieRow, Column: zombieCol}
	if !gameInstance.IsZombieOnPosition(zombieRow, zombieCol) {
		t.Errorf("Zombie should be on position [%d, %d]", zombieRow, zombieCol)
	}
}

func TestGame_GetZombieName(t *testing.T) {
	gameInstance := Game{}
	if gameInstance.GetZombieName() != "" {
		t.Errorf("Zombie name should be empty before start")
	}

	gameInstance.StartGame()
	if gameInstance.GetZombieName() == "" {
		t.Errorf("Zombie name should NOT be empty after start")
	}

	gameInstance.StopGame()
	if gameInstance.GetZombieName() != "" {
		t.Errorf("Zombie name should be empty after stop")
	}
}

func TestGame_MoveZombie(t *testing.T) {
	gameInstance := Game{}
	var position Position
	position = gameInstance.StartGame()
	lastPositionRow := position.Row

	for i := 0; i < maxColumnPosition; i++ {
		if position.Row < 0 || position.Row > maxRowPosition {
			t.Errorf("Row in zombie's position should be 0 <= x <= %d, but is %d", maxRowPosition, position.Row)
		}
		positionDelta := position.Row - lastPositionRow
		if positionDelta < -1 || positionDelta > 1 {
			t.Errorf("Zombie should be moving rows by deltas of interval [-1, 1], "+
				"row should be within distance 1 of %d, but is %d", lastPositionRow, position.Row)
		}
		if position.Column != i {
			t.Errorf("Zombie should be moving columns by 1, column should be %d but is %d", i, position.Column)
		}

		lastPositionRow = position.Row
		position = gameInstance.MoveZombie()
	}
}

func TestGame_IsGameFinished(t *testing.T) {
	gameInstance := Game{}
	if gameInstance.IsGameFinished() {
		t.Errorf("Game should not be finished before it started")
	}
	gameInstance.StartGame()
	if gameInstance.IsGameFinished() {
		t.Errorf("Game should not be finished right after it started")
	}

	for i := 0; i < maxColumnPosition-1; i++ {
		gameInstance.MoveZombie()
		if gameInstance.IsGameFinished() {
			t.Errorf("Game should be finished only after %d moves, but is finished after %d",
				maxColumnPosition, i)
		}
	}

	gameInstance.MoveZombie()
	if !gameInstance.IsGameFinished() {
		t.Errorf("Game should be finished after %d moves, but it is not",
			maxColumnPosition)
	}
}

func TestGame_StopGame(t *testing.T) {
	gameInstance := Game{}
	gameInstance.StartGame()

	gameInstance.StopGame()
	if gameInstance.gameRunning {
		t.Errorf("Game should not be running after stop")
	}
	if gameInstance.zombieName != "" {
		t.Errorf("Zombie name should be empty after stop")
	}
	if gameInstance.zombiePosition != newPosition(0, 0) {
		t.Errorf("Zombie position should be [0, 0], but is %s", gameInstance.zombiePosition.ToString())
	}
}

func TestGame_randRange(t *testing.T) {
	intervalStart := -2
	intervalEnd := 3
	for i := 0; i < 1000; i++ {
		val := randRange(intervalStart, intervalEnd)
		if val < intervalStart || val > intervalEnd {
			t.Errorf("Value should always be within inclusive interval [%d, %d], but is %d",
				intervalStart, intervalEnd, val)
		}
	}
}
