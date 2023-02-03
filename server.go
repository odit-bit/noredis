package noredis

import (
	"fmt"
	"log"
	"net"

	"github.com/odit-bit/noredis/db"
)

type Handler func(net.Conn)

type Server struct {
	Addr    string
	Handler func(net.Conn)
}

func (s *Server) ListenAndServe() error {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		fmt.Println("Failed to bind to port", s.Addr)
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("[WARNING]", err)
			continue
		}
		go s.Handler(conn)
	}

}

// will listen to addr and serve the handler
func ListenAndServe(addr string, handler Handler) error {
	s := &Server{
		Addr:    addr,
		Handler: handler,
	}
	return s.ListenAndServe()
}

//Cache will start and set noredis-cache as a handler
func Cache() Handler {
	//storage
	c := db.NewStorage()

	//command control
	cmd := NewCmd(c)

	//return as handler type
	return func(c net.Conn) {
		cmd.HandleCache(c)
	}
}
