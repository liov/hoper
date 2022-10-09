package main

import (
	"github.com/jlaffaye/ftp"
	"log"
	"time"
	"tools/backup"
)

func main() {
	c, err := ftp.Dial("192.168.137.133:2100", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	err = c.Login("anonymous", "anonymous")
	if err != nil {
		log.Fatal(err)
	}

	// Do something with the FTP conn
	Backup(c)

	if err := c.Quit(); err != nil {
		log.Fatal(err)
	}
}

func Backup(c *ftp.ServerConn) {
	backup.DCIM(c)
	backup.Pietures(c)
	//Wechat(c)
}
