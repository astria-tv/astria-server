package db

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/syndtr/goleveldb/leveldb"
	"os/user"
	"path"
	"sync"
)

// This won't survive in the longterm, we will need a RDB but for now just for playing state this will suffice
type DB struct {
	db   *leveldb.DB
	path string
}

func (self *DB) Close() {
	err := self.db.Close()
	if err == nil {
		fmt.Println("Database closed")
	} else {
		fmt.Println("Failed to close database", "err", err)
	}
}

var sharedDb struct {
	sync.Mutex
	db *DB
}

func GetSharedDB() *DB {
	sharedDb.Lock()
	defer sharedDb.Unlock()

	if sharedDb.db == nil {
		usr, err := user.Current()
		if err != nil {
			glog.Exit("Failed to determine user's home directory: ", err.Error())
		}
		sharedDb.db, err = NewDb(path.Join(usr.HomeDir, ".config", "bss", "db"))
		if err != nil {
			glog.Exit("Failed to open database: ", err.Error())
		}
	}
	return sharedDb.db
}

func NewDb(file string) (*DB, error) {
	glog.Info("Opening database at ", file)
	db, err := leveldb.OpenFile(file, nil)
	ldb := &DB{db: db, path: file}
	return ldb, err
}
