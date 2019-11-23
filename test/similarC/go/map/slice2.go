package main

import (
	"log"
	"math/rand"
	"time"
)

func main() {
	m1:= make(map[int][]int)
	for i:=0;i < 100; i++ {
		key := rand.Intn(10)
		s,ok:=m1[key]
		if !ok {
			s = make([]int,0)
		}
		s = append(s,rand.Intn(10))
	}
	log.Println(m1)

	m2:= make(map[int][]int)
	for i:=0;i < 100; i++ {
		key := rand.Intn(10)
		_,ok:=m2[key]
		if !ok {
			m2[key] = make([]int,0)
		}
		m2[key] = append(m2[key],rand.Intn(10))
	}
	log.Println(m2)

	m3:= make(map[int][]int)
	for i:=0;i < 100; i++ {
		key := rand.Intn(10)
		m3[key] = append(m3[key],rand.Intn(10))
	}
	log.Println(m3)

	log.Println(time.Now())
	data,_:=time.Now().MarshalJSON()
	log.Println(string(data))
}
