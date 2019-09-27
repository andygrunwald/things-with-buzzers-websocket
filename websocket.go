package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// WebSocketServer represents the data structure that
// keeps everything together for the web socket server
type WebSocketServer struct {
	upgrader     websocket.Upgrader
	clients      map[*websocket.Conn]bool
	buzzerStream chan buzzerHit
	logger       *logrus.Logger
}

// NewWebSocketServer will create a new web socket server instance
// "buzzer" represents the message bus of the buzzer. If a buzzer was hit,
// a message will be send into this channel
func NewWebSocketServer(buzzer chan buzzerHit, logger *logrus.Logger) *WebSocketServer {
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
		logger:       logger,
	}
	return s
}

// Broadcasting will send a single message
// to all connected websocket clients (broadcasting).
func (s *WebSocketServer) Broadcasting() {
	s.logger.Info("starting socket broadcast")
	for {
		msg := <-s.buzzerStream

		jsonMessage, _ := json.Marshal(msg)
		s.logger.WithFields(logrus.Fields{
			"msg": string(jsonMessage),
		}).Info("broadcasting message")

		// Send to every client that is currently connected
		for client := range s.clients {
			err := client.WriteMessage(websocket.TextMessage, jsonMessage)
			if err != nil {
				s.logger.WithFields(logrus.Fields{
					"err": err,
				}).Warn("websocket write message failed")
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
	s.logger.WithFields(logrus.Fields{
		"client": client.LocalAddr(),
	}).Info("new client registered")
}

// DeregisterClient signs off an existing web client to our websocket connection.
func (s *WebSocketServer) DeregisterClient(client *websocket.Conn) {
	defer s.logger.WithFields(logrus.Fields{
		"client": client.LocalAddr(),
	}).Info("client signed off")
	client.Close()
	delete(s.clients, client)
}
