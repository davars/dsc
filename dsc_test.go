package dsc

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
			out: parseResult{err: ErrMalformed, packet: Packet{}},
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
