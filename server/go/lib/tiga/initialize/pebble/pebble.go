package pebble

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/cockroachdb/pebble"
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

type PebbleDB struct {
	*pebble.DB
	Conf PebbleDBConfig
}

func (p *PebbleDB) Config() initialize.Generate {
	return &p.Conf
}

func (p *PebbleDB) SetEntity(entity interface{}) {
	if client, ok := entity.(*pebble.DB); ok {
		p.DB = client
	}
}
