package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

// webServer represents the data structure that
// keeps everything together for the webServer
// incl. the web sockets
type webServer struct {
	Listen         string
	ButtonHits     chan buttonHit
	SocketUpgrader websocket.Upgrader
	SocketClients  map[*websocket.Conn]bool
}

const (
	// buttonRed represents a physical buzzer in color red
	buttonRed string = "red"
	// buttonRed represents a physical buzzer in color green
	buttonGreen string = "green"
	// buttonRed represents a physical buzzer in color blue
	buttonBlue string = "blue"
	// buttonRed represents a physical buzzer in color yellow
	buttonYellow string = "yellow"
)

// buttonHit represents the message that will be sent
// once a button/buzzer was hit
type buttonHit struct {
	// Color is the color of the buzzer that was hit
	// see constants button* above
	Color string
}

func main() {
	log.Println("******************************************")
	log.Println("     Hardware Websocket Button Server     ")
	log.Println("******************************************")

	//
	// Start web socket and webserver
	//
	buttonHits := make(chan buttonHit, 4)
	httpServer := &webServer{
		Listen:     ":8080",
		ButtonHits: buttonHits,
		SocketUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// This is not good idea. Why?
				// See https://github.com/gorilla/websocket/issues/367
				// But we (assume to) run locally on a
				// RaspberryPi. And we want to make it work for now ;)
				return true
			},
		},
		SocketClients: make(map[*websocket.Conn]bool),
	}

	go httpServer.startWebserver()
	go httpServer.socketBroadcast()

	// Usage of https://gobot.io/for dealing with
	// physical buzzers.
	//
	// Whatever you do with the GPIO pins
	// the raw BCM2835 pinout mapping to Raspberry Pi at
	// https://godoc.org/github.com/stianeikeland/go-rpio
	// is super helpful.
	r := raspi.NewAdaptor()
	red := gpio.NewButtonDriver(r, "40")
	green := gpio.NewButtonDriver(r, "38")
	blue := gpio.NewButtonDriver(r, "36")
	yellow := gpio.NewButtonDriver(r, "32")

	work := func() {
		red.On(gpio.ButtonPush, func(data interface{}) {
			log.Println("Button red pressed")
			msg := buttonHit{
				Color: buttonRed,
			}
			buttonHits <- msg
		})

		green.On(gpio.ButtonPush, func(data interface{}) {
			log.Println("Button green pressed")
			msg := buttonHit{
				Color: buttonGreen,
			}
			buttonHits <- msg
		})

		blue.On(gpio.ButtonPush, func(data interface{}) {
			log.Println("Button blue pressed")
			msg := buttonHit{
				Color: buttonBlue,
			}
			buttonHits <- msg
		})

		yellow.On(gpio.ButtonPush, func(data interface{}) {
			log.Println("Button yellow pressed")
			msg := buttonHit{
				Color: buttonYellow,
			}
			buttonHits <- msg
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{r},
		[]gobot.Device{red, green, blue, yellow},
		work,
	)

	robot.Start()
}
