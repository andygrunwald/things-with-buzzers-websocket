package main

import (
	"log"
	"net/http"
)

// WebServer represents the data structure that
// keeps everything together for the HTTP webServer
type WebServer struct {
	listen string
	socket *WebSocketServer
	hb     *HttpBuzzer
}

// NewWebserver will create a new webserver instance
// "listen" is the ip + port combination on what the webserver should listen
// e.g., ":8080" for every interface on port 8080
// "websocket" is a WebSocketServer instance
func NewWebserver(listen string, websocket *WebSocketServer, hb *HttpBuzzer) *WebServer {
	s := &WebServer{
		listen: listen,
		socket: websocket,
		hb:     hb,
	}
	return s
}

// Start will boot up the HTTP webserver
func (s *WebServer) Start() error {
	log.Printf("Webserver starting on %s", s.listen)

	// Static file server
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// WebSocket route
	http.HandleFunc("/stream", s.websocketHandler)

	http.HandleFunc("/buzz", s.hb.buzz)

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
		log.Printf("Upgrade to websocket protocol failed: %s", err)
		return
	}

	s.socket.RegisterClient(client)
}
