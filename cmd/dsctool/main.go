package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"bitbucket.org/davars/dsc"
	"github.com/tarm/goserial"
)

func stty(args ...string) (string, error) {
	output, err := exec.Command("stty", append([]string{"-F", "/dev/tty"}, args...)...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func main() {
	device := flag.String("device", "", "IT-100 Serial Port, e.g. /dev/ttyUSB0")
	flag.Parse()

	if *device == "" {
		flag.Usage()
		return
	}

	termSettings, err := stty("-g")
	if err != nil {
		log.Fatalf("unable to save terminal settings: %v\noutput: %q", err, termSettings)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		// restore the terminal state when exiting
		<-sig
		log.Printf("exiting")
		_, err := stty(termSettings)
		if err != nil {
			log.Fatalf("failed to restore settings (%v), try running %q", err, "stty "+termSettings)
		}
		os.Exit(0)
	}()

	// disable input buffering
	stty("cbreak", "min", "1")
	// do not display entered characters on the screen
	stty("-echo")

	ser, err := serial.OpenPort(&serial.Config{Name: *device, Baud: 9600})
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		scanner := bufio.NewScanner(ser)
		for scanner.Scan() {
			bytes := scanner.Bytes()
			packet, err := dsc.Parse(bytes)
			if err != nil {
				log.Printf("error: %s\n", err)
			} else {
				log.Printf("<- %+v\n", packet)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("reading standard input: %s", err)
		}
	}()

	// initial commands
	go func() {
		// disable timestamp prefix
		send(ser, dsc.Packet{Command: 55, Data: "0"})
		// request status update
		send(ser, dsc.Packet{Command: 1})
	}()

	// Send keypresses
	go func() {
		for {
			buf := make([]byte, 1)
			n, err := os.Stdin.Read(buf)
			if n != 1 {
				panic("Read more than one byte?!")
			}
			if err != nil {
				panic(err)
			}
			// TODO: validate input
			send(ser, dsc.Packet{Command: 70, Data: string(buf)})
		}
	}()
	select {}
}

var sendMutex sync.Mutex

func send(s io.Writer, packet dsc.Packet) {
	sendMutex.Lock()
	defer sendMutex.Unlock()
	log.Printf("-> %+v\n", packet)
	_, err := s.Write(packet.Serialize())
	if err != nil {
		log.Fatalf("error sending command: %v", err)
	}
}
