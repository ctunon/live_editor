// handlers/websocket.go
package handlers

import (
	"fmt"
	"sync"

	"golang.org/x/net/websocket"
)

// WebSocketHandler manages active WebSocket connections
type WebSocketHandler struct {
	clients   map[*websocket.Conn]bool // Connected clients
	broadcast chan string              // Broadcast channel for messages
	mu        sync.Mutex               // Protects the clients map
}

// NewWebSocketHandler creates a new WebSocketHandler
func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan string),
	}
}

// HandleConnection handles new WebSocket connections
func (wsh *WebSocketHandler) HandleConnection(ws *websocket.Conn) {
	// Register new client
	wsh.mu.Lock()
	wsh.clients[ws] = true
	wsh.mu.Unlock()

	defer func() {
		// Unregister client on disconnect
		wsh.mu.Lock()
		delete(wsh.clients, ws)
		wsh.mu.Unlock()
		ws.Close()
	}()

	// Listen for incoming messages
	for {
		var message string
		if err := websocket.Message.Receive(ws, &message); err != nil {
			break
		}
		// Broadcast the message to all connected clients
		wsh.broadcast <- message
	}
}

// StartBroadcast listens for new messages and sends them to all connected clients
func (wsh *WebSocketHandler) StartBroadcast() {
	for {
		// Wait for a new message to be broadcasted
		msg := <-wsh.broadcast

		// Send the message to all connected clients
		wsh.mu.Lock()
		for client := range wsh.clients {
			err := websocket.Message.Send(client, msg)
			if err != nil {
				fmt.Println("Error sending message to client:", err)
				client.Close()
				delete(wsh.clients, client)
			}
		}
		wsh.mu.Unlock()
	}
}
