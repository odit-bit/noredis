package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"github.com/odit-bit/noredis"
	"github.com/odit-bit/noredis/db"
)

func main() {

	//read flag
	var confPath, port, password string
	flag.StringVar(&confPath, "conf", "", "config file path")
	flag.StringVar(&port, "port", "", "port number")
	flag.StringVar(&password, "password", "", "nr-server password")
	flag.Parse()

	if confPath != "" {
		config := ReadConfig(confPath)
		port = config["PORT"]
		password = config["PASSWORD"]
		config = nil
	}

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
		Addr:    port,
		Handler: cmd,
		AuthF:   authFunc,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}

}

func loggerMiddleware(logger func(a ...any), next noredis.Handler) noredis.HandlerFunc {
	return noredis.HandlerFunc(func(req *noredis.Request, res *noredis.Response) {
		logger(req.CmdName)
		next.Exe(req, res)
	})
}

func monitorMiddleWare(memStat *runtime.MemStats, after noredis.Handler) noredis.HandlerFunc {

	return noredis.HandlerFunc(func(req *noredis.Request, res *noredis.Response) {
		after.Exe(req, res)
		runtime.ReadMemStats(memStat)
		mb := memStat.Alloc / 1000
		freed := memStat.Frees
		fmt.Println("memory use:", mb, "mb")
		fmt.Println("memory freed:", freed)
	})
}
