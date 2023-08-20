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

	err = cli.Set("odit", "ganteng")
	if err != nil {
		log.Fatal(err)
	}

	result, err := cli.Get("odit")
	fmt.Println(result, err)

	cli.Set("odit", 1)
	now := time.Now()
	for i := 0; i < 1000000; i++ {
		_, err := cli.Incr("odit")
		if err != nil {
			log.Println("incr error:", err)
			break
		}
	}
	elapse := time.Since(now)
	cli.Close()
	fmt.Printf("elapsed time: %v ms \n", elapse.Milliseconds())
}
