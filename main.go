package main

import (
	"log"
	"runtime"
)

// TODO Replace log with logrus
// TODO Add software buzzer implementation
// TODO Comment everything properly
// TODO Make listen part of webserver configurable via env var

func main() {
	log.Println("******************************************")
	log.Println("      things with buzzers: websocket      ")
	log.Println("******************************************")

	// Initializing everything:
	// The websocket server, the webserver, and the buzzer implementation
	buzzerStream := make(chan buzzerHit, 4)
	websocketServer := NewWebSocketServer(buzzerStream)
	httpServer := NewWebserver(":8080", websocketServer)

	var buzzer Buzzer
	if runtime.GOARCH == "arm" {
		buzzer = NewHardwareBuzzer(buzzerStream)
		log.Println("hardware buzzer requested")
	} else {
		buzzer = NewSoftwareBuzzer(buzzerStream, ":8181")
		log.Println("software buzzer requested")
	}

	err := buzzer.Initialize()
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

	err = buzzer.Start()
	if err != nil {
		log.Fatalf("buzzer start failed: %s", err)
	}
}
