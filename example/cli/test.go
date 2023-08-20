package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/odit-bit/noredis/client"
)

func main() {
	var addr, password string

	//aggregates the args
	flag.StringVar(&addr, "addr", "8745", "address of server")
	flag.StringVar(&password, "password", "", "password server")

	//connect to noredis server
	addr = ":" + addr
	cli, err := client.Connect(addr, password)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	//stdin reader
	input := bufio.NewScanner(os.Stdin)
	input.Split(bufio.ScanLines)
	for input.Scan() {
		args := strings.Split(input.Text(), " ")
		switch args[0] {
		case "SET":
			err := cli.Set(args[1], args[2])
			if err != nil {
				log.Fatal(err)
			}

		case "GET":
			any, err := cli.Get(args[1])
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(any)

		case "INCR":
			any, err := cli.Incr(args[1])
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(any)

		case "exit":
			return

		default:
			fmt.Println("unknown command")
		}
	}

}
