package main

import (
	"github.com/sirupsen/logrus"
)

// TODO Replace log with logrus
// TODO Add software buzzer implementation
// TODO Comment everything properly
// TODO Make listen part of webserver configurable via env var

func main() {
	// logging: Configure output of timestamp
	// Otherwise it would output the time passed since beginning of execution.
	var logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
	logger.Info("******************************************")
	logger.Info("      things with buzzers: websocket      ")
	logger.Info("******************************************")

	// Initializing everything:
	// The websocket server, the webserver, and the buzzer implementation
	buzzerStream := make(chan buzzerHit, 4)
	websocketServer := NewWebSocketServer(buzzerStream, logger)
	httpServer := NewWebserver(":8080", websocketServer, logger)

	hardware := NewHardwareBuzzer(buzzerStream, logger)
	err := hardware.Initialize()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Fatal("buzzer initialisation failed")
	}

	// Start everything:
	// The websocket server, the webserver, and the buzzer implementation
	go websocketServer.Broadcasting()
	go func() {
		err := httpServer.Start()
		if err != nil {
			logger.WithFields(logrus.Fields{
				"err": err,
			}).Fatal("http server start failed")
		}
	}()

	err = hardware.Start()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Fatal("buzzer start failed")
	}
}
