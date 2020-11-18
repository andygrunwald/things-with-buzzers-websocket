package main

import (
	"flag"
	"log"
	"runtime"
)

var (
	// HTTPListenAddr represents the interface + port combination
	// where the webserver will listen on
	HTTPListenAddr = ":8080"

	// TCPListenAddr represents the interface + port combination
	// where the tcp server will listen on. The software buzzer
	// interface opens up a TCP socket to emulate buzzer
	TCPListenAddr = ":8181"

	// HardwareBuzzerSupport represents the flag to enforce
	// initialization of hardware buzzer (even on non arm architectures)
	HardwareBuzzerSupport = false

	VersionFlag = false

	// version is the software version of twb-websocket
	version = "dev"

	// commit is the current git sha of twb-websocket
	commit = "none"

	// date is the compile date of twb-websocket
	date = "unknown"
)

func main() {
	// Command line flag parsing
	flag.StringVar(&HTTPListenAddr, "http-listen-addr", LookupEnvOrString("TWB_HTTP_LISTEN_ADDR", HTTPListenAddr), "HTTP server listen address")
	flag.StringVar(&TCPListenAddr, "tcp-listen-addr", LookupEnvOrString("TWB_TCP_LISTEN_ADDR", TCPListenAddr), "TCP/Software buzzer server listen address")
	flag.BoolVar(&HardwareBuzzerSupport, "hardware-buzzer", LookupEnvOrBool("TWB_HARDWARE_BUZZER", HardwareBuzzerSupport), "Enforces initialization of hardware buzzer (even on non arm architectures)")
	flag.BoolVar(&VersionFlag, "version", VersionFlag, "Shows the current version")
	flag.Parse()

	log.Println("******************************************")
	log.Println("      things with buzzers: websocket      ")
	log.Println("******************************************")

	log.Printf("Version: %s", version)
	log.Printf("Commit: %s", commit)
	log.Printf("Date: %s", date)

	// If the version flag was hit, exit here.
	// In a normal run, we want to output the version + co information.
	if VersionFlag {
		return
	}

	// Initializing everything:
	// The websocket server, the webserver, and the buzzer implementation
	buzzerStream := make(chan buzzerHit, 4)
	websocketServer := NewWebSocketServer(buzzerStream)
	httpBuzzer := NewHttpBuzzer(buzzerStream)
	httpServer := NewWebserver(HTTPListenAddr, websocketServer, httpBuzzer)

	var buzzer Buzzer
	if runtime.GOARCH == "arm" || HardwareBuzzerSupport {
		buzzer = NewHardwareBuzzer(buzzerStream)
		log.Println("Hardware buzzer requested")
	} else {
		buzzer = NewSoftwareBuzzer(buzzerStream, TCPListenAddr)
		log.Println("Software buzzer requested")
	}

	err := buzzer.Initialize()
	if err != nil {
		log.Fatalf("Buzzer initialisation failed: %s", err)
	}

	// Start everything:
	// The websocket server, the webserver, and the buzzer implementation
	go websocketServer.Broadcasting()
	go func() {
		err := httpServer.Start()
		if err != nil {
			log.Fatalf("HTTP server start failed: %s", err)
		}
	}()

	err = buzzer.Start()
	if err != nil {
		log.Fatalf("Buzzer start failed: %s", err)
	}
}
