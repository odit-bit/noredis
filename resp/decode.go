package resp

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
)

//decode will parse resp-encoded stream(byte) into valid resp.Typ

//represent valid resp-type
type Typ struct {
	prefix byte
	body   []byte
	array  []Typ
}

// if exist it will return the value of array
func (t Typ) Array() []Typ {
	if t.prefix == ARRAY {
		return t.array
	}
	return []Typ{}
}

// convert the body field value into string
func (t Typ) String() string {
	if t.prefix == BULK || t.prefix == SIMPLE || t.prefix == INTEGER || t.prefix == ERROR {
		return string(t.body)
	}
	return ""
}

//-------------------------------------------
//Decode will parse byte into valid resp-type

//decodebulk
func decodeBulk(r *bufio.Reader) (Typ, error) {
	var t Typ
	//read length
	readLength, err := readLine(r)
	if err != nil {
		return t, err
	}
	num, err := strconv.Atoi(string(readLength))
	if err != nil {
		return t, err
	}

	// nil bulk
	if num < 1 {
		t.prefix = BULK
		t.body = []byte{}
		return t, nil
	}

	//read the body
	line, err := readLine(r)
	if err != nil {
		return t, err
	}

	t.prefix = BULK
	t.body = line
	return t, nil
}

//decodesimple
func decodeSimpleString(r *bufio.Reader) (Typ, error) {
	//"+...\r\n"
	var t Typ
	line, err := readLine(r)
	if err != nil {
		return t, err
	}

	t.prefix = SIMPLE
	t.body = line

	return t, nil
}

//ERROR
func decodeError(r *bufio.Reader) (Typ, error) {
	//"-...\r\n"
	var t Typ
	line, err := readLine(r)
	if err != nil {
		return t, err
	}
	t.prefix = ERROR
	t.body = line

	return t, nil
}

//Integer
func decodeInteger(r *bufio.Reader) (Typ, error) {
	//":...\r\n"
	var t Typ
	line, err := readLine(r)
	if err != nil {
		return t, err
	}
	t.prefix = INTEGER
	t.body = line

	return t, nil
}

//Array
func decodeArray(r *bufio.Reader) (Typ, error) {
	line, err := readLine(r)

	//fmt.Println("decodeArray", string(line))
	if err != nil {
		return Typ{}, nil
	}

	arr := []Typ{}
	num, err := strconv.Atoi(string(line))
	if err != nil {
		return Typ{}, nil
	}
	for i := 1; i <= num; i++ {
		t, err := decodeStream(r)
		if err != nil {
			return Typ{}, err
		}
		arr = append(arr, t)
		fmt.Println("arr", arr)
	}
	arrType := Typ{
		prefix: ARRAY,
		array:  arr,
	}
	return arrType, nil

}

//will parse stream of byte into appropriate  resp-type
func decodeStream(r *bufio.Reader) (Typ, error) {
	//read the first byte
	respType, err := r.ReadByte()
	if err != nil {
		return Typ{}, err
	}

	switch respType {
	case BULK:
		return decodeBulk(r)
	case SIMPLE:
		return decodeSimpleString(r)
	case INTEGER:
		return decodeInteger(r)
	case ERROR:
		return decodeError(r)
	case ARRAY:
		return decodeArray(r)
	default:
		return Typ{}, errors.New("not valid type")
	}
}

//-----------------------------------------------
//read the line until '\n'
func readLine(r *bufio.Reader) ([]byte, error) {

	var line []byte

	for {
		b, err := r.ReadByte()
		if err != nil {
			return line, err
		}

		if b == '\r' {
			continue
		}

		if b == '\n' {
			break
		}
		line = append(line, b)
	}
	return line, nil

}

func DecodeResp(r *bufio.Reader) (Typ, error) {
	return decodeStream(r)
}
