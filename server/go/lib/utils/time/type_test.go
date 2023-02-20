package timei

import (
	"encoding/json"
	"testing"
	"time"
)

type Foo struct {
	T1 UnixTime
	T2 UnixNanoTime
}

func TestType(t *testing.T) {
	foo := Foo{T1: UnixTime(time.Now()), T2: UnixNanoTime(time.Now())}
	data, _ := json.Marshal(&foo)
	t.Log(string(data))
}

func TestTimestamp(t *testing.T) {
	t.Log(time.Unix(-62135596800, 0)) // 0001-01-01 08:00:00 +0800 CST
	t.Log(time.Unix(-62135596899, 0)) // 0001-01-01 07:58:21 +0800 CST
}

type Foo1 struct {
	T1 UnionTime
	T2 UnionTime
	T3 UnionTime
	T4 UnionTime
}

func TestUnionTime(t *testing.T) {

	foo := Foo1{T1: NewUnionTime(time.Now(), 0),
		T2: NewUnionTime(time.Now(), 1),
		T3: NewUnionTime(time.Now(), 2),
		T4: NewUnionTime(time.Now(), 3),
	}
	data, _ := json.Marshal(&foo)
	t.Log(string(data)) // {"T1":"2023-02-09 15:00:49","T2":"2023-02-09","T3":1675926049,"T4":1675926049057035300}
	data = []byte(`{"T1":"2023-02-09 15:00:49","T2":"2023-02-09","T3":1675926049,"T4":1675926049057035300}`)
	foo1 := Foo1{
		T1: ZeroUnionTime(0),
		T2: ZeroUnionTime(1),
		T3: ZeroUnionTime(2),
		T4: ZeroUnionTime(3),
	}
	json.Unmarshal(data, &foo1)
	t.Log(foo1)
}

type UnionTimeInit interface {
	UnionTimeInit()
}
