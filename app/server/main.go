package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/odit-bit/noredis"
	"github.com/odit-bit/noredis/db"
)

func main() {

	if len(os.Args) != 5 {
		fmt.Println("for start using flag: ' -port port -pass password'")
		return
	}

	port := flag.String("port", "8745", "port for nr server")
	password := flag.String("pass", "", "password to connect nr server")
	// Parse command-line arguments
	flag.Parse()

	//authentication handler
	authFunc := func(pass string) bool {
		return pass == *password
	}

	// db setup
	storage := db.InitStorage()
	cmd := noredis.NewCommand(storage)

	memStat := runtime.MemStats{}
	// server
	srv := noredis.Server{
		Addr:    *port,
		Handler: monitorMiddleWare(&memStat, loggerMiddleware(log.Println, cmd)),
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
