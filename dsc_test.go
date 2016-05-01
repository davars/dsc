package dsc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type parseResult struct {
	err    error
	packet Packet
}

func TestParse(t *testing.T) {
	tests := []struct {
		in  []byte
		out parseResult
	}{
		{
			in:  []byte("6543D2"),
			out: parseResult{err: nil, packet: Packet{Command: 654, Data: "3"}},
		},
		{
			in:  []byte("6543D3"),
			out: parseResult{err: ErrChecksum, packet: Packet{}},
		},
		{
			in:  []byte("6543"),
			out: parseResult{err: ErrMalformed{in: []byte("6543")}, packet: Packet{}},
		},
		{
			in:  []byte("654332429309"),
			out: parseResult{err: nil, packet: Packet{Command: 654, Data: "3324293"}},
		},
	}

	for _, test := range tests {
		packet, err := Parse(test.in)
		assert.Equal(t, test.out, parseResult{err: err, packet: packet})
	}
}

func TestSerialize(t *testing.T) {
	tests := []struct {
		in  Packet
		out []byte
	}{
		{
			in:  Packet{Command: 654, Data: "3"},
			out: []byte("6543D2\r\n"),
		},
		{
			in:  Packet{Command: 654, Data: "3324293"},
			out: []byte("654332429309\r\n"),
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.out, test.in.Serialize())
	}
}
