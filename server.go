package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocket upgrader to upgrade HTTP connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins; restrict this in production
	},
}

// StartWebSocketServer initializes the WebSocket server and manages client connections.
func StartWebSocketServer(address string, broadcast <-chan string) {
	clients := make(map[*websocket.Conn]bool)

	// Broadcast messages to all connected clients
	go func() {
		for message := range broadcast {
			for client := range clients {
				err := client.WriteMessage(websocket.TextMessage, []byte(message))
				if err != nil {
					log.Printf("Error broadcasting to client: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}
		defer conn.Close()

		clients[conn] = true
		log.Println("New WebSocket client connected")

		// Listen for client disconnect
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Printf("WebSocket client disconnected: %v", err)
				delete(clients, conn)
				break
			}
		}
	})

	log.Printf("WebSocket server started on %s", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
