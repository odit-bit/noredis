package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/odit-bit/noredis/resp"
)

func Connect(addr, password string) (*Conn, error) {
	rwc, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	c := Conn{
		rwc:    rwc,
		reader: bufio.NewReader(rwc),
	}
	if err := c.auth(password); err != nil {
		return nil, err
	}

	return &c, nil

}

type Conn struct {
	rwc    io.ReadWriteCloser
	reader *bufio.Reader
}

func (c *Conn) auth(pass string) error {
	b := asCommand("AUTH", pass)
	_, err := c.rwc.Write(b)
	if err != nil {
		return err
	}
	_, err = resp.Unpack(c.reader)
	if err != nil {
		return err
	}
	return nil
}

func (c *Conn) Set(key string, value any) error {
	b := asCommand("SET", key, value)
	_, err := c.rwc.Write(b)
	if err != nil {
		return err
	}

	res, err := resp.Unpack(c.reader)
	if err != nil {
		return err
	}

	if res.(string) != "ok" {
		return fmt.Errorf("err: %v", res)
	}
	return nil
}

func (c *Conn) Get(key string) (any, error) {
	b := asCommand("GET", key)
	_, err := c.rwc.Write(b)
	if err != nil {
		return nil, err
	}
	return resp.Unpack(c.reader)
}

func (c *Conn) Incr(key string) (any, error) {
	b := asCommand("INCR", key)
	_, err := c.rwc.Write(b)
	if err != nil {
		return nil, err
	}
	return resp.Unpack(c.reader)
}

func (c *Conn) Close() {
	c.rwc.Close()
}

// parse args into []byte representation of noredis command
func asCommand(args ...any) []byte {
	cmd := []any{}
	cmd = append(cmd, args[0])
	for _, arg := range args[1:] {
		// num, err := strconv.Atoi(arg.(string))
		// if err != nil {
		// 	cmd = append(cmd, arg)
		// 	continue
		// }
		cmd = append(cmd, arg)
	}
	return resp.Pack(cmd)
}
