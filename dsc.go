package dsc

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type Packet struct {
	Command int
	Data    string
}

var (
	re           = regexp.MustCompile(`^((\d{3})(.*))(.{2})$`)
	ErrMalformed = errors.New("Malformed command received")
	ErrChecksum  = errors.New("Incorrect checksum")
)

func Parse(in []byte) (packet Packet, err error) {
	match := re.FindSubmatch(in)
	if len(match) != 5 {
		return Packet{}, ErrMalformed
	}

	command, err := strconv.ParseUint(string(match[2]), 10, 16)
	if err != nil {
		return Packet{}, err
	}

	expected, err := strconv.ParseUint(string(match[4]), 16, 8)
	if err != nil {
		return Packet{}, err
	}

	cksum := checksum(match[1])
	if cksum != byte(expected) {
		return Packet{}, ErrChecksum
	}

	return Packet{Command: int(command), Data: string(match[3])}, nil
}

func checksum(in []byte) (sum byte) {
	for _, b := range in {
		sum += b
	}
	return
}

func (p *Packet) Serialize() []byte {
	message := fmt.Sprintf("%03d%s", p.Command, p.Data)
	cksum := checksum([]byte(message))
	return []byte(fmt.Sprintf("%s%02X\r\n", message, cksum))
}
