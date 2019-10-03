package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

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

	b.welcomeMessage(conn)
	scanner := bufio.NewScanner(conn)
	for {
		ok := scanner.Scan()
		if !ok {
			break
		}

		b.handleMessage(scanner.Text(), conn)
	}
}

// welcomeMessage writes a welcome message to the user
// into conn
func (b *SoftwareBuzzer) welcomeMessage(conn net.Conn) {
	msg := "> Welcome to the Software Buzzer Emulator of\n"
	conn.Write([]byte(msg))
	msg = ">       things with buzzers: websocket\n"
	conn.Write([]byte(msg))
	msg = ">\n"
	conn.Write([]byte(msg))
	msg = "> Every command needs to start with \"/\"\n"
	conn.Write([]byte(msg))
	msg = "> Supported commands:\n"
	conn.Write([]byte(msg))
	msg = ">    /red\n"
	conn.Write([]byte(msg))
	msg = ">    /green\n"
	conn.Write([]byte(msg))
	msg = ">    /blue\n"
	conn.Write([]byte(msg))
	msg = ">    /yellow\n"
	conn.Write([]byte(msg))
	msg = ">    /quit\n"
	conn.Write([]byte(msg))
	msg = ">\n"
	conn.Write([]byte(msg))
	msg = "> Type your commands, and we are ready to go. Have fun!\n"
	conn.Write([]byte(msg))
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
