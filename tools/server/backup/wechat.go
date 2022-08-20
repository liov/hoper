package main

import (
	"github.com/jlaffaye/ftp"
	"log"
)

func Wechat(c *ftp.ServerConn) {
	err := Copy(c, "/Pictures/Weixin", BackUpDiskPron+"pic\\Weixin")
	if err != nil {
		log.Println(err)
	}
}
