package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

// WebServer represents the data structure that
// keeps everything together for the HTTP webServer
type WebServer struct {
	listen string
	socket *WebSocketServer
}

// NewWebserver will create a new webserver instance
// "listen" is the ip + port combination on what the webserver should listen
// e.g., ":8080" for every interface on port 8080
// "websocket" is a WebSocketServer instance
func NewWebserver(listen string, websocket *WebSocketServer) *WebServer {
	s := &WebServer{
		listen: listen,
		socket: websocket,
	}
	return s
}

// Start will boot up the HTTP webserver
func (s *WebServer) Start() error {
	log.Printf("Webserver starting on %s", s.listen)

	r := http.NewServeMux()

	// Static file server
	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/", http.StripPrefix("/static/", fs))

	// WebSocket route
	r.HandleFunc("/stream", s.websocketHandler)

	// We enable HTTP request logging here.
	// HTTP logs are right now in line with the Apache HTTPd standard format
	// This is not very beautiful, because it leads to two different log styles:
	//
	// 		2019/10/01 19:19:44 Starting socket broadcast
	//		2019/10/01 19:19:44 Webserver starting on :8080
	//		2019/10/01 19:19:44 Buzzer emulator: tcp socket starting on :8181
	//		::1 - - [01/Oct/2019:19:19:50 +0200] "GET /static/twb-jeopardy/example-season-2/game-1.json HTTP/1.1" 200 7282
	//		::1 - - [01/Oct/2019:19:19:50 +0200] "GET /favicon.ico HTTP/1.1" 404 19
	//		127.0.0.1 - - [01/Oct/2019:19:20:21 +0200] "GET /static/twb-jeopardy/seasons.json HTTP/1.1" 200 365
	//
	// As you can see, logging about the process behaviour is mixed with HTTP logging.
	// One idea would be to to adjust the HTTP logging to the normal logging structure.
	// Another idea would be to switch to logrus (see #2 https://github.com/andygrunwald/things-with-buzzers-websocket/issues/2)
	// Whatever. For now this is working and okayish.
	//
	// If you, open source contributer, read this, and you want to take action, feel free to send a PR
	// to https://github.com/andygrunwald/things-with-buzzers-websocket
	if err := http.ListenAndServe(s.listen, handlers.LoggingHandler(os.Stdout, r)); err != nil {
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
