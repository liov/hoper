package initialize

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/liov/hoper/go/v2/utils/log"
)

func P3Badger() *badger.DB{
	opts := badger.DefaultOptions("./badger")
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
