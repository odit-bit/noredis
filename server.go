package noredis

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/odit-bit/noredis/resp"
)

//noredis server implementation
//included network layer

type Handler interface {
	Exe(req *Request, res *Response)
}

type HandlerFunc func(req *Request, res *Response)

func (f HandlerFunc) Exe(req *Request, res *Response) {
	f(req, res)
}

// ============Server============

type Server struct {
	Addr        string
	Handler     Handler
	AuthF       func(string) bool
	IdleTimeout time.Duration
	rwc         []*conn
}

func (srv *Server) Serve(l net.Listener) error {

	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				break
			}
			conn := newConn(srv, c)
			srv.rwc = append(srv.rwc, conn)
			go conn.serve()
		}
	}()

	// TODO: proper shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<-sig
	err := l.Close()
	if err != nil {
		log.Println(err)
	}

	for _, c := range srv.rwc {
		c.closer.Close()
	}

	return fmt.Errorf("server is shutdown")

}

func (srv *Server) ListenAndServe() error {
	if srv.Addr == "" {
		srv.Addr = "8745"
	}
	addr := ":" + srv.Addr
	l, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("run on port", srv.Addr)
	return srv.Serve(l)
}

type serverHandler struct {
	srv *Server
}

func (sh serverHandler) Exe(req *Request, res *Response) {
	handler := sh.srv.Handler
	handler.Exe(req, res)
}

// =============Response=============

// represent requested Response
type Response struct {
	conn *conn
	req  *Request
	w    io.Writer
}

// Write write a to response writer,
// it wiil do encode before writing and
// according to type it will return error if type not allowed
func (res *Response) Write(a any) error {
	_, err := res.w.Write(resp.Pack(a))
	return err
}

//=================Conn=============

type conn struct {
	closer io.Closer
	srv    *Server
	bufR   *bufio.Reader
	w      io.Writer
}

func newConn(srv *Server, rwc io.ReadWriteCloser) *conn {
	c := conn{
		closer: rwc,
		srv:    srv,
		bufR:   bufio.NewReader(rwc),
		w:      rwc,
	}
	return &c
}

// readRequest will create response struct from connection
func (c *conn) readRequest() (*Response, error) {
	req, err := readRequest(c.bufR)
	if err != nil {
		return nil, err
	}

	res := &Response{
		conn: c,
		req:  req,
		w:    c.w,
	}

	return res, nil
}

// to authenctication the connection for the first time accepted by listener
func (c *conn) authentication() error {
	if c.srv.AuthF != nil {
		res, err := c.readRequest()
		if err != nil {
			return err
		}
		if res.req.CmdName != "AUTH" {
			err := fmt.Errorf("AUTH need authorized")
			res.Write(err)
			return err
		}

		pass := res.req.Args[0].(string)
		if ok := c.srv.AuthF(pass); !ok {
			err := fmt.Errorf("wrong password")
			res.Write(err)
			return err
		}
		res.Write("ok")
	}
	return nil
}

// serve the connection after successfully authorized
func (c *conn) serve() {
	defer func() {
		c.closer.Close()
		log.Println("conn close:")
	}()

	if err := c.authentication(); err != nil {
		return
	}

	h := serverHandler{c.srv}
	for {
		res, err := c.readRequest()
		if err != nil {
			log.Println("conn read request error:")
			break
		}
		h.Exe(res.req, res)
	}

}

//==============REQUEST============

// represent client request command

type Request struct {
	CmdName string
	Args    []any
}

// parse incoming connection as Request
func readRequest(r *bufio.Reader) (*Request, error) {
	args, err := unpackCommand(r)
	if err != nil {
		return nil, err
	}

	req := new(Request)
	req.CmdName = args[0].(string)
	req.Args = args[1:]

	return req, nil
}

// unpack command's paramater(args) from connection buffer
func unpackCommand(r io.ByteReader) ([]any, error) {
	var args []any
	unpacked, err := resp.Unpack(r)
	if err != nil {
		return nil, err
	}

	args, ok := unpacked.([]any)
	if !ok {
		return nil, fmt.Errorf("SYNTAX command should array type")
	}

	return args, err
}
