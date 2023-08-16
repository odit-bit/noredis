package resp

import (
	"io"
	"log"
	"strconv"
)

// according to redis protocol , command always in array form type
func AsCommand(r io.ByteReader, cmd *[]string) error {
	b, err := r.ReadByte()
	if err != nil {
		log.Println("parse read error:", err)
		return err
	}
	if b != '*' {
		log.Println("not array form")
		return err
	}
	size, err := getSize(r)
	if err != nil {
		log.Println("array length not recognized")
		return err
	}

	//iteration by size

	for i := 0; i < size; i++ {
		//skip read the prefix
		_, _ = r.ReadByte()
		ele, err := readCrlf(r)
		if err != nil {
			log.Println("parse error read crlf:", err)
		}

		*cmd = append(*cmd, string(ele))
	}
	return nil
}

// will read until crlf
func readCrlf(br io.ByteReader) ([]byte, error) {
	var buf []byte
	for {
		b, err := br.ReadByte()
		if err != nil {
			return nil, err
		}

		switch b {
		case '\r':
			continue
		case '\n':
			return buf, nil
		default:
			buf = append(buf, b)
		}
	}
}

func getSize(br io.ByteReader) (int, error) {
	arr, err := readCrlf(br)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(arr))
}

//=============================================
//=============================================

// func AsReply(a any) []byte {
// 	var reply []byte
// 	fmt.Printf("as reply :%v, type :%t \n", a, a)
// 	switch t := a.(type) {
// 	case noredis.SIMPLE:
// 		VALUE := string(t)
// 		reply = append(reply, []byte(_SIMPLE+VALUE+_CRLF)...)
// 	case noredis.NULL:
// 		reply = append(reply, []byte("_"+_CRLF)...)

// 	case noredis.BLOB:
// 		VALUE := string(t)
// 		SIZE := strconv.Itoa(len(VALUE))
// 		reply = append(reply, []byte(_BLOB+SIZE+_CRLF+VALUE+_CRLF)...)

// 	case error:
// 		VALUE := t.Error()
// 		reply = append(reply, []byte(_ERROR+VALUE+_CRLF)...)

// default:
//
//		return nil //fmt.Errorf("conn write error: unknown type")
//	}
//
// fmt.Println("reply:", string(reply))
// return reply //nil
// }
