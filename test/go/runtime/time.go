package main

import (
	"fmt"
	"log"
	"time"
)
//windows的Sleep，2ms算可以了？
func main()  {
	//这三行没啥意义，证明分时系统不精确而已
	var StartTime = time.Now().Unix()
	time.Sleep(1 * time.Second)
	fmt.Println(time.Now().Unix() - StartTime)

	now:=time.Now().UnixNano()
	start:=now
	for{
		t:=0
		for start >=now{
			t++
			time.Sleep(1)
			now = time.Now().UnixNano()
		}
		start = now
		if t > 1 {
			log.Printf("t:%d",t)
		}
	}
}
