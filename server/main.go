package main

import (
	"fmt"

	"github.com/odit-bit/noredis"
)

func main() {

	// start cache engine
	fmt.Println("start no-redis-cache engine")
	noredis.ListenAndServe("0.0.0.0:6379", noredis.Cache())

}
