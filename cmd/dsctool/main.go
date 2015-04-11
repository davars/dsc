package main

import (
	"bitbucket.org/davars/dsc"
	"bufio"
	"flag"
	"fmt"
	"github.com/tarm/goserial"
	"io"
	"log"
	"os"
	"os/exec"
)

func main() {
	device := flag.String("device", "", "IT-100 Serial Port, e.g. /dev/ttyUSB0")
	flag.Parse()

	if *device == "" {
		flag.Usage()
		return
	}

	c := &serial.Config{Name: *device, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	wait := make(chan bool)
	go func() {
		scanner := bufio.NewScanner(s)
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

	// Request status update
	go func() {
		send(s, dsc.Packet{Command: 1})
	}()

	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// restore the echoing state when exiting
	fmt.Fprintln(os.Stderr, "WARNING: echo disabled, 'stty -F /dev/tty echo' to restore it!") // if the next line fails
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()

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
			send(s, dsc.Packet{Command: 70, Data: string(buf)})
		}
	}()
	<-wait
}

func send(s io.Writer, packet dsc.Packet) {
	log.Printf("-> %+v\n", packet)
	s.Write(packet.Serialize())
}
