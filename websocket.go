package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocketServer represents the data structure that
// keeps everything together for the web socket server
type WebSocketServer struct {
	upgrader     websocket.Upgrader
	clients      map[*websocket.Conn]bool
	buzzerStream chan buzzerHit
}

// NewWebSocketServer will create a new web socket server instance
// "buzzer" represents the message bus of the buzzer. If a buzzer was hit,
// a message will be send into this channel
func NewWebSocketServer(buzzer chan buzzerHit) *WebSocketServer {
	s := &WebSocketServer{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// This is not good idea. Why?
				// See https://github.com/gorilla/websocket/issues/367
				// But we (assume to) run locally on a
				// Raspberry Pi. And we want to make it work for now ;)
				//
				// If you read this and have a good idea on how to fix
				// this, be brave, change it and send us a PR to
				// https://github.com/andygrunwald/things-with-buzzers-websocket
				return true
			},
		},
		clients:      make(map[*websocket.Conn]bool),
		buzzerStream: buzzer,
	}
	return s
}

// Broadcasting will send a single message
// to all connected websocket clients (broadcasting).
func (s *WebSocketServer) Broadcasting() {
	log.Println("Starting socket broadcast ...")
	for {
		msg := <-s.buzzerStream

		jsonMessage, _ := json.Marshal(msg)
		log.Printf("Broadcasting message: %v\n", string(jsonMessage))

		// Send to every client that is currently connected
		for client := range s.clients {
			err := client.WriteMessage(websocket.TextMessage, jsonMessage)
			if err != nil {
				log.Printf("Websocket error: %s", err)
				s.DeregisterClient(client)
			}
		}
	}
}

// UpgradeConnection upgrades the HTTP connection to the WebSocket protocol.
// TODO: Deprecated: Use websocket.Upgrader instead.
func (s *WebSocketServer) UpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	client, err := s.upgrader.Upgrade(w, r, nil)
	return client, err
}

// RegisterClient registers a new web client to our websocket connection.
func (s *WebSocketServer) RegisterClient(client *websocket.Conn) {
	s.clients[client] = true
}

// DeregisterClient signs off an existing web client to our websocket connection.
func (s *WebSocketServer) DeregisterClient(client *websocket.Conn) {
	client.Close()
	delete(s.clients, client)
}
