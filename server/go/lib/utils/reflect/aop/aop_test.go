package aop

import (
	"log"
	"testing"
)

var foo1 = func() {
	log.Println("foo1")
}

func llog() { log.Println("log") }

func TestAop(t *testing.T) {
	Invoke(llog, &foo1)
	foo1()
	foo1()
}
