package ws

import (
	"zombie-game/src/chTypes"
	"zombie-game/src/game"
)

type Hub struct {

	// Channel for receiving Message broadcast
	broadcast chan chTypes.Broadcast

	// Storage for games hosted by different players
	gameInstances map[string]*GameData

	// Channel that gets notified when game should end, receives hosting player's name
	stopGame chan string
}

type GameData struct {
	// Name of the player that is hosting the game
	playerName string

	gameInstance *game.Game
	// Receivers that can register to receive broadcast messages about this game
	clientMessageReceivers []chan string
	// Channel used to make new shot attempts
	shotAttempt chan chTypes.ShotAttempt
}

func NewHub() *Hub {
	return &Hub{
		broadcast:     make(chan chTypes.Broadcast),
		gameInstances: make(map[string]*GameData),
		stopGame:      make(chan string),
	}
}

func NewGameData(playerName string) *GameData {
	return &GameData{
		playerName:             playerName,
		gameInstance:           &game.Game{HostingPlayerName: playerName},
		clientMessageReceivers: make([]chan string, 0),
		shotAttempt:            make(chan chTypes.ShotAttempt),
	}
}

func (hub *Hub) CreateNewGame(hostingPlayerName string) {
	hub.gameInstances[hostingPlayerName] = NewGameData(hostingPlayerName)
	hub.gameInstances[hostingPlayerName].gameInstance.Run(hub.broadcast,
		hub.gameInstances[hostingPlayerName].shotAttempt,
		hub.stopGame)
}

func (gameData *GameData) RegisterBroadcastReceiver(receiver chan string) {
	gameData.clientMessageReceivers = append(gameData.clientMessageReceivers, receiver)
}

func (h *Hub) Run() {
	for {
		select {
		case hostingPlayer := <-h.stopGame:
			delete(h.gameInstances, hostingPlayer)

		case message := <-h.broadcast:
			for _, element := range h.gameInstances[message.HostingPlayerName].clientMessageReceivers {
				element <- message.Message
			}
		}
	}
}
