package db

// storage is in-memory storage of key-value (map) store data

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	DEFAULT_CLEAN_INTERVAL = time.Duration(5 * time.Minute)
)

type database struct {
	stop chan struct{}
	wg   sync.WaitGroup
	mu   sync.RWMutex

	data map[key]Data
}

//will call for every interval to clean-up expired data
func (db *database) cleanUp(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		log.Println("db clean-up")
		select {
		case <-db.stop:
			return
		case <-t.C:
			db.mu.Lock()
			for key, value := range db.data {
				if value.expireAt.Unix() <= time.Now().Unix() {
					delete(db.data, key)
				}
			}
			db.mu.Unlock()
		}
	}
}

//---------------------------------------------------

//set the key to value and return the len of value
//if the key exist it will update the old value
func (db *database) Set(k key, v Data) int64 {
	db.mu.Lock()
	defer db.mu.Unlock()
	if v.Expired != 0 {
		v.setPX()
	}
	db.data[k] = v
	return 0
}

//get data by key
func (db *database) Get(k key) (string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	data, ok := db.data[k]
	fmt.Println("db.Get", ok)
	if !ok {
		return string(data.Value), errors.New("no such key")
	}
	if data.isExpired() {
		fmt.Println("data expired")
		delete(db.data, k)
		return "", nil //errors.New("is expired")
	}
	fmt.Println("data available")
	return data.Value, nil
}

//Delete existed data
func (db *database) Delete(k key) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	_, ok := db.data[k]
	if !ok {
		return errors.New("no such key")
	}
	delete(db.data, k)
	return nil
}

//will close the stop channel and wait all go-routine to finish
func (db *database) Stop() {
	close(db.stop)
	db.wg.Wait()
}

//run the database clean-up goroutines
func (db *database) CleanUp(interval time.Duration) {
	db.wg.Add(1)
	go func() {
		defer db.wg.Done()
		db.cleanUp(interval)
	}()
}

//create new db instance
func New() *database {
	return &database{
		stop: make(chan struct{}),
		wg:   sync.WaitGroup{},
		mu:   sync.RWMutex{},
		data: map[key]Data{},
	}
}

type Storage struct {
	data *database
}

func (s *Storage) Data() *database {
	return s.data
}

func NewStorage() *Storage {
	db := New()
	return &Storage{
		data: db,
	}
}
