package main

import (
	"bitbucket.org/davars/dsc"
	"bufio"
	"github.com/tarm/goserial"
	"log"
)

func main() {
	c := &serial.Config{Name: "/dev/cu.usbserial", Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(s)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		packet, err := dsc.Parse(bytes)
		if err != nil {
			log.Printf("error: %s\n", err)
		} else {
			log.Printf("%+v\n", packet)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("reading standard input: %s", err)
	}
}
