package main

import (
	"log"
	"testing"

	"github.com/odit-bit/noredis"
	"github.com/odit-bit/noredis/db"
)

func Benchmark_Incr(b *testing.B) {
	password := ""
	//authentication handler
	authFunc := func(pass string) bool {
		return pass == password
	}

	// db setup
	storage := db.InitStorage()
	cmd := noredis.NewCommand(storage)

	// memStat := runtime.MemStats{}
	// server
	srv := noredis.Server{
		Addr:    "8745",
		Handler: cmd,
		AuthF:   authFunc,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
