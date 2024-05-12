package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func read1(path string) {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	buf := make([]byte, 1024)
	for {
		n, err := fi.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if 0 == n {
			break
		}
	}
}

func read3(path string) {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	_, err = ioutil.ReadAll(fi)
}

func main() {

	file := "./text.txt" //找一个大的文件，如日志文件
	start := time.Now()
	read1(file)
	t1 := time.Now()
	fmt.Printf("Cost time %v\n", t1.Sub(start))
	read3(file)
	t3 := time.Now()
	fmt.Printf("Cost time %v\n", t3.Sub(t1))

	sr := strings.NewReader("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	buf := bufio.NewReaderSize(sr, 0) //默认16
	b := make([]byte, 10)

	fmt.Println("==", buf.Buffered()) // 0
	s, _ := buf.Peek(5)
	fmt.Printf("%d ==  %q\n", buf.Buffered(), s) //
	nn, er := buf.Discard(3)
	fmt.Println(nn, er)

	for n, err := 0, error(nil); err == nil; {
		fmt.Printf("Buffered:%d ==Size:%d== n:%d==  b[:n] %q ==  err:%v\n", buf.Buffered(), buf.Size(), n, b[:n], err)
		n, err = buf.Read(b)
		fmt.Printf("Buffered:%d ==Size:%d== n:%d==  b[:n] %q ==  err: %v == s: %s\n", buf.Buffered(), buf.Size(), n, b[:n], err, s)
	}

	fmt.Printf("%d ==  %q\n", buf.Buffered(), s)
}
