package pebble

import (
	"github.com/liov/hoper/server/go/lib/initialize"
	"github.com/liov/hoper/server/go/lib/utils/log"
	"github.com/cockroachdb/pebble"
)

type Config struct {
	DirName string
}

func (conf *Config) Build() *pebble.DB {
	if conf.DirName == "" {
		return nil
	}
	db, err := pebble.Open(conf.DirName, nil)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func (conf *Config) Generate() interface{} {
	return conf.Build()
}

type DB struct {
	*pebble.DB
	Conf Config
}

func (p *DB) Config() initialize.Generate {
	return &p.Conf
}

func (p *DB) SetEntity(entity interface{}) {
	if client, ok := entity.(*pebble.DB); ok {
		p.DB = client
	}
}

func (p *DB) Close() error {
	if p.DB == nil {
		return nil
	}
	return p.DB.Close()
}
