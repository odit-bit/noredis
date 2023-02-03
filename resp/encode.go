package resp

import (
	"fmt"
	"strconv"
	"strings"
)

// A Redis server replies to clients, sending any valid RESP data type as a reply.
// encode will parse a reply or response into the stream(byte)

//encode toBULK-string type
func EncodeString(s string) []byte {

	f := fmt.Sprintf("%c%v\r\n%v\r\n", BULK, len(s), s)
	return []byte(f)
}

//encode to SimpleString
func EncodeSimple(s string) []byte {
	return []byte(string(SIMPLE) + s + crlf)
}

//encode to integer-type
func encodeInt(i int) []byte {
	f := fmt.Sprintf(":%v\r\n", i)
	return []byte(f)
}

// wrapper of encodeInt
func EncodeInt(i int) []byte {
	return encodeInt(i)
}

// encode to error-type
func EncodeErr(err error) []byte {
	f := fmt.Sprintf("%cERROR %v\r\n", ERROR, err.Error())
	return []byte(f)
}

// will choose Encode(Type) function to use
func typeEncoder(str string) []byte {
	//check if the string is digit
	num, err := strconv.Atoi(str)
	if err == nil {
		//is digit
		return encodeInt(num)
	}
	//is not digit
	return EncodeString(str)
}

//encode string command to RESP array-type
func Encode(s string) []byte {
	a := strings.Split(s, " ")

	//array format
	format := []byte(fmt.Sprintf("*%d\r\n", len(a)))

	for _, v := range a {
		b := typeEncoder(v)
		//format += fmt.Sprintf("%s", b)
		format = append(format, b...)
	}
	return format
}

func EncodeBulk(s string) []byte {
	num := len(s)
	if num < 1 {
		format := []byte("$-1\r\n") //nil bulk
		return format
	}
	return []byte(fmt.Sprintf("$%d\r\n%v\r\n", num, s))
}
