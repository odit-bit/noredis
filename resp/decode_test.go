package resp

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_decodeSimple(t *testing.T) {
	r := require.New(t)
	input := []byte("+Ping\r\n")
	buf := bytes.NewBuffer(input)
	reader := bufio.NewReader(buf)

	//read the first byte '+'
	_, err := reader.ReadByte()
	r.NoError(err)

	expect := Typ{
		prefix: SIMPLE,
		body:   []byte("Ping"),
	}

	actual, err := decodeSimpleString(reader)
	r.NoError(err)
	r.Equal(expect, actual)
}

func Test_decodeStream(t *testing.T) {
	r := require.New(t)

	input := []byte(
		"*4\r\n" +
			"$5\r\nHELLO\r\n" +
			"$5\r\nWORLD\r\n" +
			"$-1\r\n" +
			":10\r\n")
	buf := bytes.NewBuffer(input)
	reader := bufio.NewReader(buf)

	expect := Typ{
		prefix: ARRAY,
		body:   []byte{},
		array: []Typ{
			{
				prefix: BULK,
				body:   []byte("HELLO"),
			},
			{
				prefix: BULK,
				body:   []byte("WORLD"),
			},
			{
				prefix: BULK,
				body:   []byte{},
			},
			{
				prefix: INTEGER,
				body:   []byte("10"),
			},
		},
	}
	value, err := decodeStream(reader)
	r.NoError(err)
	r.Equal(expect.array, value.array, "should the same")

}
