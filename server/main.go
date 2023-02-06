package main

import (
	"fmt"

	"github.com/odit-bit/noredis"
)

func main() {

	// start cache engine
	fmt.Println("start no-redis-cache engine")
	_, err := noredis.ListenAndServe("localhost:6379", noredis.Cache())
	if err != nil {
		fmt.Println(err)
	}

}
