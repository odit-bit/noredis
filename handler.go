package noredis

import (
	"fmt"
)

// // handler is wired all the commponent
type Command struct {
	storage Storage
}

func NewCommand(db Storage) *Command {
	// http.Server
	cmd := &Command{
		// conn:    conn,
		storage: db,
	}
	return cmd
}

func (cmd *Command) Exe(req *Request, res *Response) {
	var result any
	var err error

	switch req.CmdName {
	case "SET":
		// cmd.storage.Add(req.Args[0].(string), req.Args[1])
		result, err = cmd.set(req, res)

	case "GET":
		any, ok := cmd.storage.LookupRead(req.Args[0].(string))
		if !ok {
			err = fmt.Errorf("key not exist")
		}
		result = any

	case "INCR":
		result, err = cmd.incr(req.Args[0].(string))

	default:
		err = fmt.Errorf("illegal Command")
	}

	if err != nil {
		res.Write(err)
		return
	}

	res.Write(result)
}

func (cmd *Command) set(req *Request, res *Response) (any, error) {
	args := req.Args
	var opt Setoptions
	if len(args) > 2 {
		err := opt.fromArgs(args[2:])
		if err != nil {
			return nil, err
		}
	}
	return Set(cmd.storage, args[0].(string), args[1], &opt)
}

func (cmd *Command) incr(key string) (int, error) {
	v, ok := cmd.storage.LookupRead(key)
	if !ok {
		num := 1
		cmd.storage.Add(key, num)
		return num, nil
	}

	switch num := v.(type) {
	case int:
		num += 1
		cmd.storage.Add(key, num)
		return num, nil
	default:
		return 0, fmt.Errorf("INCR value is not number")
	}
}
