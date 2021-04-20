package initialize

import (
	"github.com/cockroachdb/pebble"
	"github.com/liov/hoper/go/v2/utils/log"
)

type PebbleDBConfig struct {
	DirName string
}

func (conf *PebbleDBConfig) generate() *pebble.DB {
	db, err := pebble.Open(conf.DirName, nil)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func (conf *PebbleDBConfig) Generate() interface{} {
	return conf.generate()
}
