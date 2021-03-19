package initialize

import (
	"github.com/cockroachdb/pebble"
	"github.com/liov/hoper/go/v2/utils/log"
)

func P3Pebble() *pebble.DB{
	db,err:=pebble.Open("./pebble",nil)
	if err != nil {
		log.Fatal(err)
	}
	return db
}