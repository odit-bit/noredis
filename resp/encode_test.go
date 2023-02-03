package resp_test

import (
	"errors"
	"testing"

	"github.com/odit-bit/noredis/resp"
	"github.com/stretchr/testify/require"
)

func TestEncodeString(t *testing.T) {
	r := require.New(t)

	/* encode bulk-string */
	expectedByte := []byte("$3\r\nKEY\r\n")
	v := resp.EncodeString("KEY")
	r.Equal(expectedByte, v, "byte should equal")

	expectedByte = []byte("$0\r\n\r\n")
	v = resp.EncodeString("")
	r.Equal(expectedByte, v, "byte should equal")

	/* encode int */
	expectedByte = []byte(":0\r\n")
	num := resp.EncodeInt(0)
	r.Equal(expectedByte, num, "integer should same")

	expectedByte = []byte(":1505\r\n")
	num = resp.EncodeInt(1505)
	r.Equal(expectedByte, num, "integer should same")

	/* encode Error */
	expectedByte = []byte("-ERROR Message\r\n")
	err := resp.EncodeErr(errors.New("Message"))
	r.Equal(expectedByte, err, "error should has message")

	/* encode Array */
	//for client to send command, only use array consisted of a bulk-string type
	expectedByte = []byte("*2\r\n$3\r\nGET\r\n$7\r\nCOMMAND\r\n")
	cmd := resp.Encode("GET COMMAND")
	r.Equal(expectedByte, cmd, "Command should the same")
}
