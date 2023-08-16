package resp

import "strconv"

func EncodeSimple(s string) []byte {
	VALUE := string(s)
	return append([]byte{}, []byte(_SIMPLE+VALUE+_CRLF)...)
}

func EncodeNumber(s int) []byte {
	VALUE := strconv.Itoa(s)
	return append([]byte{}, []byte(_NUMBER+VALUE+_CRLF)...)
}

func EncodeBLOB(s []byte) []byte {
	VALUE := s
	SIZE := strconv.Itoa(len(VALUE))
	res := append([]byte{}, []byte(_BLOB+SIZE+_CRLF)...)
	res = append(res, VALUE...)
	res = append(res, []byte(_CRLF)...)
	return res
}

func EncodeError(err error) []byte {
	VALUE := err.Error()
	return append([]byte{}, []byte(_ERROR+VALUE+_CRLF)...)
}

func EncodeArray(agg []any) []byte {
	size := len(agg)
	arr := []byte(_AGG_ARRAY + strconv.Itoa(size) + _CRLF)

	for _, ele := range agg {
		arr = append(arr, Pack(ele)...)
	}
	return arr
}
