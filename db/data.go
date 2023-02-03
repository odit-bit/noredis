package db

//in-mem key-value store  datastructure (type)

import (
	"fmt"
	"time"
)

type key interface{}

//represent memory stored-data
type Data struct {
	Value    string
	Expired  int
	expireAt time.Time
}

//set data expiration
func (d *Data) setPX() {
	d.expireAt = time.Now().Add(time.Duration(d.Expired) * time.Millisecond)
	fmt.Println("Set with", d.expireAt)
}

//check if data is expired
func (d *Data) isExpired() bool {
	if d.Expired == 0 {
		return false
	}
	return d.expireAt.Before(time.Now())
}
