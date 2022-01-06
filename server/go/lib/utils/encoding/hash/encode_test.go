package hash

import (
	"reflect"
	"testing"
	"time"

	"github.com/actliboy/hoper/server/go/lib/utils/log"
)

type Foo struct {
	Time time.Time
	Bar
}

type Bar struct {
	Int int
}

func TestMarshal(t *testing.T) {
	e := new(encodeState)
	u := &Foo{Time: time.Now(), Bar: Bar{Int: 1}}
	e.encode("", reflect.ValueOf(u))
	for i := 0; i < len(e.strings); i += 2 {
		log.Info(e.strings[i], e.strings[i+1])
	}
}
