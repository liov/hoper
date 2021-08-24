package inject_dao

import (
	"github.com/cockroachdb/pebble"
	"github.com/liov/hoper/server/go/lib/utils/log"
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
