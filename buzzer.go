package main

import (
	"log"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

const (
	// buzzerColorRed represents a physical buzzer in color red
	buzzerColorRed string = "red"
	// buttonRed represents a physical buzzer in color green
	buzzerColorGreen string = "green"
	// buttonRed represents a physical buzzer in color blue
	buzzerColorBlue string = "blue"
	// buttonRed represents a physical buzzer in color yellow
	buzzerColorYellow string = "yellow"

	// buzzerPinRed represents the Raspberry Pi GPIO pin where the
	// physical buzzer in color red is connnected
	buzzerPinRed = "40"

	// buzzerPinGreen represents the Raspberry Pi GPIO pin where the
	// physical buzzer in color green is connnected
	buzzerPinGreen = "38"

	// buzzerPinBlue represents the Raspberry Pi GPIO pin where the
	// physical buzzer in color blue is connnected
	buzzerPinBlue = "36"

	// buzzerPinYellow represents the Raspberry Pi GPIO pin where the
	// physical buzzer in color yellow is connnected
	buzzerPinYellow = "32"
)

// buzzerHit represents the message that will
// be sent once a buzzer was hit.
type buzzerHit struct {
	// Color is the color of the buzzer that was hit
	// see constants buzzerColor*
	Color string
}

// Buzzer represents an interface of different kind of buzzers.
// Typical usecases are hardware buzzers or software buzzers
// e.g., to emulate the hardware.
type Buzzer interface {
	// Initialize sets up everything that is needed to run the buzzer.
	Initialize() error
	// Start boots up the buzzer and ensures that they are ready to hit.
	Start() error
}

// HardwareBuzzer is the implementation of the Buzzer
// interface for physical hardware buzzers.
// See https://github.com/andygrunwald/things-with-buzzers-hardware
// for more details
type HardwareBuzzer struct {
	buzzerStream chan buzzerHit
	robot        *gobot.Robot
}

// NewHardwareBuzzer returns a new instance of
// the hardware buzzer stream
func NewHardwareBuzzer(buzzer chan buzzerHit) Buzzer {
	b := &HardwareBuzzer{
		buzzerStream: buzzer,
	}
	return b
}

// Initialize sets up everything that is needed to run the buzzer.
func (b *HardwareBuzzer) Initialize() error {
	// Usage of https://gobot.io/ for
	// dealing with physical buzzers.
	//
	// Whatever you do with the GPIO pins
	// the raw BCM2835 pinout mapping to Raspberry Pi at
	// https://godoc.org/github.com/stianeikeland/go-rpio
	// is super helpful.
	r := raspi.NewAdaptor()
	red := gpio.NewButtonDriver(r, buzzerPinRed)
	green := gpio.NewButtonDriver(r, buzzerPinGreen)
	blue := gpio.NewButtonDriver(r, buzzerPinBlue)
	yellow := gpio.NewButtonDriver(r, buzzerPinYellow)

	work := func() {
		red.On(gpio.ButtonPush, func(data interface{}) {
			log.Println("Button red pressed")
			msg := buzzerHit{
				Color: buzzerColorRed,
			}
			b.buzzerStream <- msg
		})

		green.On(gpio.ButtonPush, func(data interface{}) {
			log.Println("Button green pressed")
			msg := buzzerHit{
				Color: buzzerColorGreen,
			}
			b.buzzerStream <- msg
		})

		blue.On(gpio.ButtonPush, func(data interface{}) {
			log.Println("Button blue pressed")
			msg := buzzerHit{
				Color: buzzerColorBlue,
			}
			b.buzzerStream <- msg
		})

		yellow.On(gpio.ButtonPush, func(data interface{}) {
			log.Println("Button yellow pressed")
			msg := buzzerHit{
				Color: buzzerColorYellow,
			}
			b.buzzerStream <- msg
		})
	}

	b.robot = gobot.NewRobot("buttonBot",
		[]gobot.Connection{r},
		[]gobot.Device{red, green, blue, yellow},
		work,
	)
	return nil
}

// Start boots up the buzzer and ensures that they are ready to hit.
func (b *HardwareBuzzer) Start() error {
	err := b.robot.Start()
	return err
}
