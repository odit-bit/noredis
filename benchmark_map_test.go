package noredis

import (
	"bytes"
	"testing"
)

// benchmark for retrieve map[string]string/*string/[]byte

func Benchmark_pointer_map(b *testing.B) {
	key := "name"
	value := string(bytes.Repeat([]byte{255}, 8096))
	bucket := map[string]*string{key: &value}
	var buf bytes.Buffer

	for i := 0; i < b.N; i++ {
		buf.Reset()
		v := bucket[key]
		buf.WriteString(*v)
	}
}

// var real string

func Benchmark_value_map(b *testing.B) {
	key := "name"
	value := string(bytes.Repeat([]byte{255}, 8096))
	bucket := map[string]string{key: value}
	var buf bytes.Buffer

	for i := 0; i < b.N; i++ {
		buf.Reset()
		v := bucket[key]
		buf.WriteString(v)
	}
}

func Benchmark_bytes_map(b *testing.B) {
	key := "name"
	value := bytes.Repeat([]byte{255}, 8096)
	bucket := map[string][]byte{key: value}
	var buf bytes.Buffer

	var v []byte
	for i := 0; i < b.N; i++ {
		buf.Reset()
		v = bucket[key]
		buf.Write(v)

	}

}

func Benchmark_Alloc(b *testing.B) {
	type req struct {
		cmd  string
		args []any
	}

	rq := func() *req {
		r := req{
			cmd:  "SET",
			args: []any{"as", 1, 2, "dwdwd"},
		}
		return &r
	}

	var request req

	for i := 0; i < b.N; i++ {
		result := rq()
		request = *result
	}

	_ = request
}
