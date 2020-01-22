package ws

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"zombie-game/src/game"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	playerName   string
	gameInstance *game.Game
	// Channel for receiving zombie position updates
	zombiePosition chan game.Position
	// Channel for stopping a running game
	gameStop chan struct{}
}

func newClient(hub *Hub, conn *websocket.Conn) Client {
	return Client{hub: hub, conn: conn, send: make(chan []byte, 256), zombiePosition: make(chan game.Position)}
}

// readHandler received messages from the websocket connection and handles them using handleMessage.
// One goroutine is started for each connection.
// readHandler can also write to the connection, but it only writes responses to events received
//  from the client, not system events
func (c *Client) readHandler() {
	defer func() {
		c.hub.unregister <- c
		//close(c.gameStop)
		c.conn.Close()
	}()
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		message := string(messageBytes)
		log.Println("Received message:", message)

		c.handleMessage(message)
	}
}

// writeHandler sends events to the websocket connection.
// One goroutine is started for each connection.
// All system events (not received from client) to one specific connection go through this handler,
//  but not necessarily all writes (see readHandler).
func (c *Client) writeHandler() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case position := <-c.zombiePosition:
			output := fmt.Sprintf("WALK %s %d %d", c.gameInstance.GetZombieName(), position.Row, position.Column)
			c.writeString(output)

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the client
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(hub, conn)
	client.hub.register <- &client

	// start read and write handlers
	go client.writeHandler()
	go client.readHandler()
}
