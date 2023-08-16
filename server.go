package noredis

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/odit-bit/noredis/resp"
)

type Handler interface {
	Exe(req *Request, res *Response)
}

type HandlerFunc func(req *Request, res *Response)

func (f HandlerFunc) Exe(req *Request, res *Response) {
	f(req, res)
}

// ============Server============

type Server struct {
	Addr    string
	Handler Handler
}

func (srv *Server) Serve(l net.Listener) error {
	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		conn := newConn(srv, c)
		conn.serve()
	}
}

func (srv *Server) ListenAndServe() error {
	if srv.Addr == "" {
		srv.Addr = ":8745"
	}
	l, err := net.Listen("tcp", srv.Addr)

	if err != nil {
		log.Fatal(err)
	}

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

type Response struct {
	conn *conn
	req  *Request
	w    io.Writer
}

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

func (c *conn) serve() {
	defer func() {
		err := c.closer.Close()
		log.Println("conn close", err)
	}()
	for {
		res, err := c.readRequest()
		if err != nil {
			log.Println(err)
			break
		}
		serverHandler{c.srv}.Exe(res.req, res)
	}

}

//==============REQUEST============

// represent client request command

type Request struct {
	CmdName string
	Args    []any
}

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
