package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"

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

	// res, err := cmd.Set("22a10", "hello world", 10000)
	res, err := cmd.Connect("")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("auth:", res)

	res, err = cmd.Set("odit", 1, 0)
	fmt.Println("res", res, err)

	for {
		res, err = cmd.Get("odit")
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Println(res)
		time.Sleep(10 * time.Millisecond)
	}
}

// ============ implement command
// repesent Command package
type Command struct {
	conn *bufio.ReadWriter
	buf  []byte
}

func (c *Command) Connect(pass string) (any, error) {
	arr := []any{"AUTH", pass}
	c.conn.Write(resp.Pack(arr))
	c.conn.Flush()
	return resp.Unpack(c.conn)
}

func (c *Command) Set(key string, value any, expire int) (any, error) {
	// it will make arg as []any
	arr := []any{"SET", key, value}
	// if Pack invocked with rest parameter it will return only the first argument
	c.conn.Write(resp.Pack(arr))
	c.conn.Flush()
	// read the response
	return resp.Unpack(c.conn)
}

func (c *Command) Get(key string) (any, error) {
	// it will make arg as []any
	arr := []any{"GET", key}
	// if Pack invocked with rest parameter it will return only the first argument
	c.conn.Write(resp.Pack(arr))
	c.conn.Flush()
	// read the response
	return resp.Unpack(c.conn)
}
