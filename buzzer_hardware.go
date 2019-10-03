// +build !windows

package main

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

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
			sendBuzzerHit(b.buzzerStream, buzzerColorRed)
		})

		green.On(gpio.ButtonPush, func(data interface{}) {
			sendBuzzerHit(b.buzzerStream, buzzerColorGreen)
		})

		blue.On(gpio.ButtonPush, func(data interface{}) {
			sendBuzzerHit(b.buzzerStream, buzzerColorBlue)
		})

		yellow.On(gpio.ButtonPush, func(data interface{}) {
			sendBuzzerHit(b.buzzerStream, buzzerColorYellow)
		})
	}

	b.robot = gobot.NewRobot("buzzerBot",
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
