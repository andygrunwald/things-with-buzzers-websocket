package main

import (
	"log"
)

const (
	// buzzerColorRed represents a physical buzzer in color red
	buzzerColorRed string = "red"
	// buzzerColorGreen represents a physical buzzer in color green
	buzzerColorGreen string = "green"
	// buzzerColorBlue represents a physical buzzer in color blue
	buzzerColorBlue string = "blue"
	// buzzerColorYellow represents a physical buzzer in color yellow
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

// sendBuzzerHit adds a new buzzer hit message (color c)
// into the buzzer stream b
func sendBuzzerHit(b chan buzzerHit, c string) {
	log.Printf("Buzzer pressed: %s", c)
	msg := buzzerHit{
		Color: c,
	}
	b <- msg
}
