package main

import (
	"fmt"
	"log"
	"time"

	"github.com/odit-bit/noredis/client"
)

func main() {
	addr := ":8745"
	cli, err := client.Connect(addr, "")
	if err != nil {
		log.Fatal(err)
	}

	err = cli.Set("22a10", "hellow")
	if err != nil {
		log.Fatal(err)
	}

	result, err := cli.Get("22a10")
	fmt.Println(result, err)

	cli.Set("22a10", 1)
	now := time.Now()
	for i := 0; i < 1000; i++ {
		_, err := cli.Incr("22a10")
		if err != nil {
			log.Println("incr error:", err)
			break
		}
	}
	elapse := time.Since(now)
	cli.Close()
	fmt.Printf("elapsed time: %v ms \n", elapse.Milliseconds())
}
