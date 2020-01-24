package ws

import (
	"fmt"
	"strconv"
	"strings"
	"zombie-game/src/chTypes"
)

// game commands
const (
	startToken             = "START"
	joinToken              = "JOIN"
	shootToken             = "SHOOT"
	startMessageTokenCount = 2
	joinMessageTokenCount  = 3
	shootMessageTokenCount = 3
)

func (client *Client) handleMessage(message string) {
	tokens := strings.Split(message, " ")

	switch tokens[0] {
	case startToken:
		if len(tokens) != startMessageTokenCount {
			client.writeString(fmt.Sprintf("%s command must have %s tokens", startToken, startMessageTokenCount))
			break
		}
		client.handleGameStartCommand(tokens[1])
	case joinToken:
		if len(tokens) != joinMessageTokenCount {
			client.writeString(fmt.Sprintf("%s command must have %s tokens", joinToken, joinMessageTokenCount))
			break
		}
		client.handleJoinGameCommand(tokens[1], tokens[2])
	case shootToken:
		if len(tokens) != shootMessageTokenCount {
			client.writeString(fmt.Sprintf("%s command must have %s tokens", shootToken, shootMessageTokenCount))
			break
		}
		client.handleShotCommand(tokens[1], tokens[2])
	default:
		client.writeString("Command not recognized")
	}
}

func (client *Client) handleGameStartCommand(playerName string) {
	if _, ok := client.hub.gameInstances[playerName]; ok {
		client.writeString(fmt.Sprintf("Game is already hosted by %s, join with 'JOIN %s <yourName>'",
			playerName, playerName))
		return
	}
	client.hub.CreateNewGame(playerName)

	client.hostingPlayerName = playerName
	client.yourName = playerName
	client.joinedGame = client.hub.gameInstances[playerName]
	client.joinedGame.RegisterBroadcastReceiver(client.broadcastReceiver)

	gameStartMessage := fmt.Sprintf("Game started by %s", playerName)
	client.broadcastMessage(gameStartMessage)
}

func (client *Client) handleJoinGameCommand(hostingPlayerName string, yourName string) {
	if _, ok := client.hub.gameInstances[hostingPlayerName]; !ok {
		client.writeString(fmt.Sprintf("Game hosted by %s doesn't exist, start one with 'START <yourName>'",
			hostingPlayerName))
		return
	}

	client.hostingPlayerName = hostingPlayerName
	client.yourName = yourName
	client.joinedGame = client.hub.gameInstances[hostingPlayerName]
	client.joinedGame.RegisterBroadcastReceiver(client.broadcastReceiver)

	gameJoinMessage := fmt.Sprintf("%s joined game hosted by %s", yourName, hostingPlayerName)
	client.broadcastMessage(gameJoinMessage)
}

func (client *Client) handleShotCommand(targetRow string, targetColumn string) {
	if client.joinedGame == nil {
		client.writeString("You need to start or join a game first: 'START <player>' / 'JOIN <player> <yourName>'")
		return
	}

	row, err := strconv.Atoi(targetRow)
	if err != nil {
		client.writeString(fmt.Sprintf("Cannot parse token: %s", targetRow))
		return
	}
	column, err := strconv.Atoi(targetColumn)
	if err != nil {
		client.writeString(fmt.Sprintf("Cannot parse token: %s", targetColumn))
		return
	}

	shotAttempt := chTypes.ShotAttempt{
		PlayerName: client.yourName,
		Row:        row,
		Column:     column,
	}
	client.joinedGame.shotAttempt <- shotAttempt
}
