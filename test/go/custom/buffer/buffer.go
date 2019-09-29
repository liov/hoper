package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"log"
)

func main() {
	buffer:=bytes.Buffer{}
	err:=binary.Write(&buffer,binary.BigEndian,int64(123))
	if err!=nil{
		log.Println(err)
	}
	s:=hex.Dump(buffer.Bytes())
	log.Println(s)
}
