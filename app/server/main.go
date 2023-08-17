package main

import (
	"log"

	"github.com/odit-bit/noredis"
	"github.com/odit-bit/noredis/db"
)

func main() {

	storage := db.InitStorage()
	cmd := noredis.NewCommand(storage)

	// server
	srv := noredis.Server{
		Addr:    ":8745",
		Handler: loggerMiddleware(log.Println, cmd),
		AuthF:   Auth,
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

func Auth(pass string) bool {
	return pass == "hire_me!!"
}
