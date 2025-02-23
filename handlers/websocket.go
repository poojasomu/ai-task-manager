package handlers

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // Connected clients
var broadcast = make(chan string) // Message channel
var mutex = sync.Mutex{} // Protects clients map

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}

// WebSocket handler
func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}

		fmt.Println("Received WebSocket Message:", string(msg))

		// Broadcast message to all connected clients
		broadcast <- string(msg)

		// Echo back to sender
		conn.WriteMessage(messageType, msg)
	}
}

// Broadcast task updates to all clients
func BroadcastTaskUpdate(message string) {
	mutex.Lock()
	defer mutex.Unlock()
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Println("WebSocket send failed:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

// Background goroutine to listen for updates
func StartBroadcast() {
	for {
		msg := <-broadcast
		BroadcastTaskUpdate(msg)
	}
}
