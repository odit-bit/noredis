package db

import "github.com/odit-bit/noredis"

type Storage interface {
	Add(key string, value any)
	LookupRead(key string) (any, bool)
}

//==================DB================

// db implement Storage interface{}
type db struct {
	m map[string]any
}

func InitStorage() noredis.Storage {
	db := &db{
		m: map[string]any{},
	}
	return db
}

func (db *db) Add(key string, value any) {
	db.m[key] = value
}

func (db *db) LookupRead(key string) (any, bool) {
	v, ok := db.m[key]
	return v, ok
}
