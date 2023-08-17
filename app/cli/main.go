package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/odit-bit/noredis/resp"
)

func main() {
	addr := "103-181-183-201.nevacloud.io:8745"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		addr := ":8745"
		conn, err = net.Dial("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}
	}

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	buf := make([]byte, 1024)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		args := scanner.Text()
		cmd := strings.Split(args, " ")
		rw.Writer.Write(AsCommand(cmd))
		rw.Flush()

		n, _ := rw.Reader.Read(buf[0:])
		fmt.Println(">>", string(buf[:n]))
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
