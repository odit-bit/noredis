package noredis

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/odit-bit/noredis/db"
	"github.com/odit-bit/noredis/resp"
)

//Command will handling resp.typ into appropriate command

type Command struct {
	cache *db.Storage
	//conn  net.Conn
}

func (cmd *Command) HandleCache(conn net.Conn) {
	cache := cmd.cache
	defer conn.Close()
	for {
		//DecodeResp return valid resp-type
		respType, err := resp.DecodeResp(bufio.NewReader(conn))
		if err != nil {
			_, err := conn.Write(resp.EncodeErr(err))
			fmt.Println("[DEBUG]DecodeResp", err)
			break
		}
		//process command
		cmdArg := respType.Array()[0].String()
		args := respType.Array()[1:]
		fmt.Println("[DEBUG]parse cmd", cmd, args)

		switch cmdArg {
		case "ping", "PING":
			//fmt.Println("switch to PING", cmd)
			n, err := conn.Write(pingCommand(args))
			fmt.Println("server write", err, "byte", n)

		case "ECHO", "echo":
			_, err := conn.Write(echoCommand(args))
			fmt.Println("server write", err)

		case "set", "SET":
			result, err := cmd.Set(args)
			if err != nil {
				fmt.Println("[DEBUG-Set]", err)
			}
			conn.Write(result)
		case "get", "GET":
			v, err := cache.Data().Get(args[0].String())
			if err != nil {
				conn.Write(resp.EncodeErr(err))
				continue
			}
			conn.Write(resp.EncodeBulk(v))

		default:
			conn.Write([]byte("-ERR unknown command '" + cmdArg + "'\r\n"))

		}
	}
}

func (cmd *Command) Set(args []resp.Typ) ([]byte, error) {
	if len(args) < 2 {
		//conn.Write(resp.EncodeErr(fmt.Errorf("need more argument , got %d", len(args))))
		return resp.EncodeErr(fmt.Errorf("need more argument , got %d", len(args))), nil
	}
	data := db.Data{
		Value:   args[1].String(),
		Expired: 0,
	}
	//with expired
	if len(args) > 2 {
		if args[2].String() == "px" {
			fmt.Println("with px", args[3].String())
			ms, err := strconv.Atoi(args[3].String())
			if err != nil {
				resp.EncodeErr(fmt.Errorf("ERR px value (%v) is not integer \r\n %v", ms, err.Error()))
			}
			data.Expired = ms
		}
	}

	n := cmd.cache.Data().Set(args[0].String(), data)
	_ = n
	return resp.Encode("OK"), nil
}

//Ping handler
func pingCommand(args []resp.Typ) []byte {
	if len(args) == 0 {
		response := resp.EncodeSimple("PONG")
		fmt.Println(string(response))
		return response
	}
	response := resp.EncodeSimple("PONG " + args[0].String())
	fmt.Println(string(response))
	return response
}

//Echo command
func echoCommand(args []resp.Typ) []byte {
	//no argument
	if len(args) == 0 {
		return resp.EncodeErr(errors.New("no args got 0"))
	}
	return []byte(fmt.Sprintf("$%d\r\n%v\r\n", len(args[0].String()), args[0].String()))
}

//------------------------------------

func NewCmd(cache *db.Storage) *Command {
	return &Command{
		cache: cache,
	}
}
