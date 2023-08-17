package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/odit-bit/noredis/resp"
)

func main() {
	addr := ":8745"

	conn, err := net.Dial("tcp", addr)

	if err != nil {
		if err != nil {
			log.Fatal(err)
		}
	}

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	cmd := Command{
		conn: rw,
	}

	res, err := cmd.Set("22a10", "hello world", 10000)
	fmt.Println(res)

}

// ============ implement command
// repesent Command package
type Command struct {
	conn *bufio.ReadWriter
	buf  []byte
}

func (c *Command) Set(key string, value any, expire int) (any, error) {
	// it will make arg as []any
	arr := []any{"SET", key, value, "PX", expire}
	// if Pack invocked with rest parameter it will return only the first argument
	c.conn.Write(resp.Pack(arr))
	c.conn.Flush()
	// read the response
	return resp.Unpack(c.conn)
}

func (c *Command) Get(key string) []byte {
	return resp.Pack([]any{"GET", key})
}
