package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// startWebserver will start the webserver incl.
// websocket - Nothing more.
func (server *webServer) startWebserver() {
	log.Printf("Starting websocket on %s ...\n", server.Listen)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/socket", server.websocketHandler)
	if err := http.ListenAndServe(server.Listen, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// websocketHandler registers new websocket clients
func (server *webServer) websocketHandler(w http.ResponseWriter, r *http.Request) {
	c, err := server.SocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	// Register new client
	server.SocketClients[c] = true
}

// socketBroadcast will send a single message
// to all connected websocket clients (broadcasting).
func (server *webServer) socketBroadcast() {
	log.Println("Starting socket broadcast ...")
	for {
		msg := <-server.ButtonHits

		jsonMessage, _ := json.Marshal(msg)
		log.Printf("Broadcasting message: %v\n", string(jsonMessage))

		// Send to every client that is currently connected
		for client := range server.SocketClients {
			err := client.WriteMessage(websocket.TextMessage, jsonMessage)
			if err != nil {
				log.Printf("Websocket error: %s", err)
				client.Close()
				delete(server.SocketClients, client)
			}
		}
	}
}
