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

	// sig := make(chan os.Signal, 1)
	// signal.Notify(sig, os.Interrupt)

	// for {
	// 	select {
	// 	case <-sig:
	// 		fmt.Println("receive signal")
	// 	case <-time.After(1 * time.Second):
	// 		fmt.Println("Hello in a loop")
	// 		// default:
	// 		// 	//noredis.ListenAndServe("0.0.0.0:6379", noredis.Cache())
	// 		// 	fmt.Println("dafault")
	// 	}
	// }
}
