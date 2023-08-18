package resp

import (
	"fmt"
	"io"
	"strconv"
)

//encode/decode resp protocol

const _CRLF string = "\r\n"
const _BLOB string = "$"
const _SIMPLE string = "+"
const _ERROR string = "-"
const _NUMBER string = ":"

const _AGG_ARRAY string = "*"

// ===========ENCODER================

func Pack(a ...any) []byte {
	var ele []byte
	for _, arg := range a {
		switch t := arg.(type) {
		case string:
			ele = EncodeSimple(t)
		case []byte:
			ele = EncodeBLOB(t)
		case int:
			ele = EncodeNumber(t)
		case error:
			ele = EncodeError(t)
		case []any:
			ele = EncodeArray(t)
		default:
			panic("unknown type")
		}
	}
	return ele
}

// ===========DECODER================

func Unpack(r io.ByteReader) (any, error) {
	var ele any
	var err error

	t, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	switch t {
	case '+':
		ele, err = DecodeSimple(r)
	case ':':
		ele, err = DecodeNumber(r)
	case '$':
		ele, err = DecodeBLOB(r)
	case '*':
		ele, err = DecodeArray(r)
	default:
		err = fmt.Errorf("illegal type")
	}

	if err != nil {
		return nil, err
	}
	return ele, nil
}

func DecodeSimple(r io.ByteReader) (string, error) {
	b, err := readCrlf(r)
	return string(b), err
}

func DecodeBLOB(r io.ByteReader) ([]byte, error) {
	blobSize, err := getSize(r)
	if err != nil {
		return nil, err
	}

	blob, err := readCrlf(r)
	if err != nil {
		return nil, err
	}
	if blobSize != len(blob) {
		return nil, fmt.Errorf("length is difference")
	}

	return blob, nil
}

func DecodeNumber(r io.ByteReader) (int, error) {
	b, err := readCrlf(r)
	if err != nil {
		return 0, err
	}

	num, err := strconv.Atoi(string(b))
	if err != nil {
		return 0, fmt.Errorf("NUMBER value is not number's representation")
	}

	return num, nil

}

func DecodeArray(r io.ByteReader) ([]any, error) {
	size, err := getSize(r)
	if err != nil {
		return nil, err
	}

	arr := make([]any, 0, size)
	for i := 0; i < size; i++ {
		ele, err := Unpack(r)
		if err != nil {
			return nil, err
		}
		arr = append(arr, ele)
	}
	return arr, nil
}
