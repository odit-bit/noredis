package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/odit-bit/noredis/resp"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s host:port \n", os.Args[0])
	}
	service := os.Args[1]

	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkErr(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkErr(err)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		//get input from user
		scanner.Scan()
		text := scanner.Text()
		if text == "//quit" {
			break
		}

		//write resp-encoded message
		encodeMsg := resp.Encode(text)
		_, err = conn.Write(encodeMsg)
		checkErr(err)

		reader := bufio.NewReader(conn)
		fmt.Println(readResponse(reader))
	}

	//close the connection
	conn.Close()

}

func readResponse(reader *bufio.Reader) (string, error) {
	response, err := resp.DecodeResp(reader)
	if err != nil {
		return "", err
	}

	str := ""
	if len(response.Array()) < 1 {
		return response.String(), nil
	}

	for _, v := range response.Array() {
		str += v.String() + " "
		fmt.Println("read Response ", v.String())
	}

	return str, nil
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
