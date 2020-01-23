package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"strconv"
	"strings"
	"zombie-game/src/game"
)

// game commands
const (
	startToken             = "START"
	shootToken             = "SHOOT"
	startMessageTokenCount = 2
	shootMessageTokenCount = 3
)

func (client *Client) handleMessage(message string) {
	tokens := strings.Split(message, " ")

	switch tokens[0] {
	case startToken:
		if len(tokens) != startMessageTokenCount {
			client.writeString("START command must have 2 tokens")
			break
		}
		client.handleGameStartCommand(tokens[1])
	case shootToken:
		if len(tokens) != shootMessageTokenCount {
			client.writeString("SHOOT command must have 3 tokens")
			break
		}
		client.handleShotCommand(tokens[1], tokens[2])
	default:
		client.writeString("Command not recognized")
	}
}

func (client *Client) handleGameStartCommand(playerName string) {
	client.playerName = playerName
	gameInstance := game.Game{}
	client.gameInstance = &gameInstance
	client.gameStop = client.gameInstance.Run(client.zombiePosition, client.gameLost)
}

func (client *Client) handleShotCommand(targetRow string, targetColumn string) {
	if client.gameInstance == nil {
		client.writeString("You need to start the game first: 'START <player>'")
		return
	}

	x, err := strconv.Atoi(targetRow)
	if err != nil {
		client.writeString(fmt.Sprintf("Cannot parse token: %s", targetRow))
		return
	}
	y, err := strconv.Atoi(targetColumn)
	if err != nil {
		client.writeString(fmt.Sprintf("Cannot parse token: %s", targetColumn))
		return
	}

	if client.gameInstance.IsZombieOnPosition(x, y) {
		output := fmt.Sprintf("BOOM %s wins!", client.playerName)
		close(client.gameStop)
		client.writeString(output)
	} else {
		client.writeString("MISS")
	}
}

func (c *Client) writeString(s string) {
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}
	w.Write([]byte(s))
	w.Close()
}
