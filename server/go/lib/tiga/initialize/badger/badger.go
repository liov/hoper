package badger

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/dgraph-io/badger/v3"
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

type Consumer struct {
	*badger.DB
	Conf BadgerDBConfig
}

func (b *Consumer) Config() initialize.Generate {
	return &b.Conf
}

func (b *Consumer) SetEntity(entity interface{}) {
	if client, ok := entity.(*badger.DB); ok {
		b.DB = client
	}
}
