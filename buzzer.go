package main

import (
	"bufio"
	"log"
	"net"
	"os"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
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
			log.Println("Buzzer red pressed")
			msg := buzzerHit{
				Color: buzzerColorRed,
			}
			b.buzzerStream <- msg
		})

		green.On(gpio.ButtonPush, func(data interface{}) {
			log.Println("Buzzer green pressed")
			msg := buzzerHit{
				Color: buzzerColorGreen,
			}
			b.buzzerStream <- msg
		})

		blue.On(gpio.ButtonPush, func(data interface{}) {
			log.Println("Buzzer blue pressed")
			msg := buzzerHit{
				Color: buzzerColorBlue,
			}
			b.buzzerStream <- msg
		})

		yellow.On(gpio.ButtonPush, func(data interface{}) {
			log.Println("Buzzer yellow pressed")
			msg := buzzerHit{
				Color: buzzerColorYellow,
			}
			b.buzzerStream <- msg
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

// SoftwareBuzzer is the implementation of the Buzzer
// interface for "software" buzzers.
// We use this interface, if we don't have the hardware
// buzzers ready or at the same location. By this,
// we are still able to play with it and create new frontends.
// See https://github.com/andygrunwald/things-with-buzzers-hardware
// for more details
type SoftwareBuzzer struct {
	buzzerStream chan buzzerHit
	listen       string
}

// NewSoftwareBuzzer returns a new instance of
// the software buzzer stream
func NewSoftwareBuzzer(buzzer chan buzzerHit, listen string) Buzzer {
	b := &SoftwareBuzzer{
		buzzerStream: buzzer,
		listen:       listen,
	}
	return b
}

// Initialize sets up everything that is needed to run the buzzer.
func (b *SoftwareBuzzer) Initialize() error {
	return nil
}

// Start boots up the buzzer and ensures that they are ready to hit.
func (b *SoftwareBuzzer) Start() error {
	listener, err := net.Listen("tcp", b.listen)
	if err != nil {
		return err
	}
	defer listener.Close()
	log.Printf("Buzzer emulator: tcp socket starting on %s", b.listen)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Buzzer emulator: tcp accept failed: %s", err)
		} else {
			go b.handleConnection(conn)
		}
	}
}

// handleConnection takes care about a single connection.
// A connection mostly send several messages.
func (b *SoftwareBuzzer) handleConnection(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	log.Printf("Buzzer emulator: client connected: %s", remoteAddr)
	defer log.Printf("Buzzer emulator: client disconnected: %s", remoteAddr)

	scanner := bufio.NewScanner(conn)
	for {
		ok := scanner.Scan()
		if !ok {
			break
		}

		b.handleMessage(scanner.Text(), conn)
	}
}

// handleMessage takes care about every single message from conn.
func (b *SoftwareBuzzer) handleMessage(message string, conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	log.Printf("Buzzer emulator: message received from client %s: %s", remoteAddr, message)

	if len(message) > 0 && message[0] == '/' {
		switch {
		case message == "/red":
			sendBuzzerHit(b.buzzerStream, buzzerColorRed)
			b.sendOKMessage(conn)
		case message == "/green":
			sendBuzzerHit(b.buzzerStream, buzzerColorGreen)
			b.sendOKMessage(conn)
		case message == "/blue":
			sendBuzzerHit(b.buzzerStream, buzzerColorBlue)
			b.sendOKMessage(conn)
		case message == "/yellow":
			sendBuzzerHit(b.buzzerStream, buzzerColorYellow)
			b.sendOKMessage(conn)

		case message == "/quit":
			log.Println("Buzzer emulator: i was told to shut down")
			conn.Write([]byte("I'm shutting down now.\n"))
			log.Println("Buzzer emulator: bye bye")
			os.Exit(0)

		default:
			conn.Write([]byte("Unrecognized command.\n"))
		}
	}
}

// sendOKMessage sends a "OK" message to the client conn.
func (b *SoftwareBuzzer) sendOKMessage(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	resp := "OK\n"
	log.Printf("Buzzer emulator: message sent to client %s: %s", remoteAddr, resp)
	conn.Write([]byte(resp))
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
