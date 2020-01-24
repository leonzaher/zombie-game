package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
	"zombie-game/src/chTypes"
)

const (
	// Time allowed to write a Message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong Message from the peer.
	pongWait = 1 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Game that you joined
	joinedGame *GameData
	// Name of hosting player
	hostingPlayerName string
	// Name of your player
	yourName string
	// Channel for receiving broadcasting messages
	broadcastReceiver chan string
}

func newClient(hub *Hub, conn *websocket.Conn) Client {
	return Client{hub: hub,
		conn:              conn,
		broadcastReceiver: make(chan string)}
}

// readHandler received messages from the websocket connection and handles them using handleMessage.
// One goroutine is started for each connection.
// readHandler can also write to the connection, but it only writes responses to events received
//  from the client, not system events
func (c *Client) readHandler() {
	defer func() {
		log.Println(fmt.Sprintf("Player %s disconnected", c.yourName))
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
		log.Println("Received Message:", message)

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
		case message := <-c.broadcastReceiver:
			c.writeString(message)

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) broadcastMessage(message string) {
	c.hub.broadcast <- chTypes.Broadcast{HostingPlayerName: c.hostingPlayerName, Message: message}
}

func (c *Client) writeString(s string) {
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}
	w.Write([]byte(s))
	w.Close()
}

// serveWs handles websocket requests from the client
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(hub, conn)

	// start read and write handlers
	go client.writeHandler()
	go client.readHandler()
}
