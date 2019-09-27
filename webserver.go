package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// WebServer represents the data structure that
// keeps everything together for the HTTP webServer
type WebServer struct {
	listen string
	socket *WebSocketServer
	logger *logrus.Logger
}

// NewWebserver will create a new webserver instance
// "listen" is the ip + port combination on what the webserver should listen
// e.g., ":8080" for every interface on port 8080
// "websocket" is a WebSocketServer instance
func NewWebserver(listen string, websocket *WebSocketServer, logger *logrus.Logger) *WebServer {
	s := &WebServer{
		listen: listen,
		socket: websocket,
		logger: logger,
	}
	return s
}

// Start will boot up the HTTP webserver
func (s *WebServer) Start() error {
	s.logger.WithFields(logrus.Fields{
		"listen": s.listen,
	}).Info("webserver starting")

	// Static file server
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// WebSocket route
	http.HandleFunc("/stream", s.websocketHandler)

	if err := http.ListenAndServe(s.listen, nil); err != nil {
		return err
	}
	return nil
}

// websocketHandler handles the HTTP route for new websocket clients.
// It upgrades the connection and registers the client.
func (s *WebServer) websocketHandler(w http.ResponseWriter, r *http.Request) {
	client, err := s.socket.UpgradeConnection(w, r)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"err": err,
		}).Warn("upgrade to websocket protocol failed")
		return
	}

	s.socket.RegisterClient(client)
}
