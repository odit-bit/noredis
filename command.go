package noredis

import (
	"errors"
)

type Storage interface {
	Add(key string, value any)
	LookupRead(key string) (any, bool)
}

var ErrBadSyntax = errors.New("SYNTAX SET key value   [ex|px value] | [nx|xx] ... ")
var ErrExpireNotNumber = errors.New("EX/PX value should string represntation of number")
var ErrKeyShouldExist = errors.New("XX key not exist")
var ErrKeyShouldNotExist = errors.New("NX key is exist")
var ErrUnknownCommand = errors.New("unknown command")
var ErrValueNotNumber = errors.New("INCR value of key not number")

func Set(db Storage, key string, value any, opt *Setoptions) (any, error) {

	_, ok := db.LookupRead(key)

	if opt.keyShouldExist {
		if !ok {
			return nil, ErrKeyShouldExist
		}
	}
	if opt.keyShouldNotExist {
		if ok {
			return nil, ErrKeyShouldNotExist
		}
	}

	db.Add(key, value)
	// if opt.hasExpire {
	// 	dur, err := strconv.Atoi(opt.expireValue)
	// 	if err != nil {
	// 		return ErrExpireNotNumber
	// 	}

	// 	var unix int64
	// 	if opt.expireInSecond {
	// 		unix = time.Now().Add(time.Duration(dur) * time.Second).Unix()
	// 	} else {
	// 		unix = time.Now().Add(time.Duration(dur) * time.Millisecond).Unix()
	// 	}
	// 	err = db.AddExpire(key, unix)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return "ok", nil
}

type Setoptions struct {
	keyShouldExist    bool
	keyShouldNotExist bool
	hasExpire         bool
	expireInSecond    bool
	expireValue       string
}

func (opt *Setoptions) fromArgs(args []any) error {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "nx":
			opt.keyShouldNotExist = true
			continue
		case "xx":
			opt.keyShouldExist = true
			continue

		case "px":
			opt.expireInSecond = false
			if i+1 < len(args) {
				opt.expireValue = args[i+1].(string)
				opt.hasExpire = true
				i++
				continue
			}

		case "ex":
			opt.expireInSecond = true
			if i+1 < len(args) {
				opt.expireValue = args[i+1].(string)
				opt.hasExpire = true
				i++
				continue
			}

		default:
			return ErrBadSyntax
		}
	}
	return nil
}
