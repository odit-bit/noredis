package noredis

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

	"github.com/odit-bit/noredis/db"
)

type Handler func(net.Conn)

type Server struct {
	Addr    string
	Handler func(net.Conn)
}

func (s *Server) ListenAndServe() error {

	//start listen to tcp
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		fmt.Println("Failed to bind to port", s.Addr)
		return err
	}

	//accept the connection
	sig := make(chan os.Signal, 1)
	wg := sync.WaitGroup{}
	signal.Notify(sig, os.Interrupt)
	count := 0

	for {

		select {
		case <-sig:
			fmt.Println("receive signal")
			listener.Close()
			wg.Wait()
			fmt.Println("finish process")
		default:
			fmt.Println("wait incoming connection...")
			conn, err := listener.Accept() //block
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					fmt.Println("listener closed")
					return nil
				}
				log.Println("[WARNING]", err)
				continue
			}

			//
			fmt.Println("--------------incoming connection: ", count+1)
			wg.Add(1)
			go func() {
				s.Handler(conn)
				wg.Done()
			}()
			count++
		}
	}

}

//------------------------------------------------------------------------------

// will listen to addr and serve the handler
func ListenAndServe(addr string, handler Handler) (*Server, error) {
	s := &Server{
		Addr:    addr,
		Handler: handler,
	}
	//add wg to wait the s.quit chan
	return s, s.ListenAndServe()
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
