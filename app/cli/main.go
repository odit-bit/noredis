package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/odit-bit/noredis/resp"
)

func main() {
	env := os.Args[1]
	var addr string
	switch env {
	case "local":
		addr = ":8745"
	case "remote":
		addr = "103-181-183-201.nevacloud.io:8745"
	default:
		addr = env
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connect to address :", addr)

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	buf := make([]byte, 1024)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	done := make(chan struct{})
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// read Op
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				n, err := rw.Reader.Read(buf[0:])
				if err != nil {
					// log.Println(err)
					return
				}
				fmt.Println(">>", string(buf[:n]))
			}
			// break
		}
	}()

	//write op
	for scanner.Scan() {
		select {
		case <-sig:
			close(done)
			conn.Close()
			return
		default:
			args := scanner.Text()
			cmd := strings.Split(args, " ")
			rw.Writer.Write(AsCommand(cmd))
			rw.Flush()
		}
	}

}

func AsCommand(args []string) []byte {
	// var cmd []byte
	// cmd = append(cmd, '*')
	// cmd = append(cmd, []byte(strconv.Itoa(len(args)))...)
	// cmd = append(cmd, "\r\n"...)

	// for _, arg := range args {
	// 	cmd = append(cmd, []byte("+"+arg+"\r\n")...)
	// }

	// return cmd
	cmd := []any{}
	cmd = append(cmd, args[0])
	for _, arg := range args[1:] {
		num, err := strconv.Atoi(arg)
		if err != nil {
			cmd = append(cmd, arg)
			continue
		}
		cmd = append(cmd, num)
	}
	return resp.Pack(cmd)
}
