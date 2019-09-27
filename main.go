package main

import (
	"log"
)

// TODO Replace log with logrus
// TODO Add software buzzer implementation
// TODO Comment everything properly
// TODO Make listen part of webserver configurable via env var

func main() {
	log.Println("******************************************")
	log.Println("     Hardware Websocket Buzzer Server     ")
	log.Println("******************************************")

	// Initializing everything:
	// The websocket server, the webserver, and the buzzer implementation
	buzzerStream := make(chan buzzerHit, 4)
	websocketServer := NewWebSocketServer(buzzerStream)
	httpServer := NewWebserver(":8080", websocketServer)

	hardware := NewHardwareBuzzer(buzzerStream)
	err := hardware.Initialize()
	if err != nil {
		log.Fatalf("buzzer initialisation failed: %s", err)
	}

	// Start everything:
	// The websocket server, the webserver, and the buzzer implementation
	go websocketServer.Broadcasting()
	go func() {
		err := httpServer.Start()
		if err != nil {
			log.Fatalf("http server start failed: %s", err)
		}
	}()

	err = hardware.Start()
	if err != nil {
		log.Fatalf("buzzer start failed: %s", err)
	}
}
