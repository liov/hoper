package badger

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/liov/hoper/server/go/lib/utils/log"
)

type BadgerDBConfig struct {
	Path string
}

func (conf *BadgerDBConfig) Build() *badger.DB {
	opts := badger.DefaultOptions(conf.Path)
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

type Consumer struct {
	*badger.DB
	Conf BadgerDBConfig
}

func (c *Consumer) Config() any {
	return &c.Conf
}

func (c *Consumer) SetEntity() {
	c.DB = c.Conf.Build()
}

func (c *Consumer) Close() error {
	return c.DB.Close()
}
