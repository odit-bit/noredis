package main

import (
	"log"

	"github.com/odit-bit/noredis"
	"github.com/odit-bit/noredis/db"
)

func main() {

	//when server recieve incoming connection from socket (tcp)
	// it will create conn object that cann read(parse) the command from client and write(reply)to that connection

	storage := db.InitStorage()
	cmd := noredis.NewCommand(storage)

	srv := noredis.Server{
		Addr:    ":8745",
		Handler: loggerMiddleware(log.Println, cmd),
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
