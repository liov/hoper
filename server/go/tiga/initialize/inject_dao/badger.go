package inject_dao

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/liov/hoper/v2/utils/log"
)

type BadgerDBConfig struct {
	Path string
}

func (conf *BadgerDBConfig) generate() *badger.DB {
	opts := badger.DefaultOptions(conf.Path)
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func (conf *BadgerDBConfig) Generate() interface{} {
	return conf.generate()
}
