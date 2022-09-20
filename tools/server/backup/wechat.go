package backup

import (
	"github.com/jlaffaye/ftp"
	"log"
)

func Wechat(c *ftp.ServerConn) {
	err := Copy(c, "/Pictures/Weixin", BackUpDiskPron+"pic\\Weixin", true)
	if err != nil {
		log.Println(err)
	}
	err = Copy(c, "/DCIM/1034735436", BackUpDiskPron+"pic\\1034735436", true)
	if err != nil {
		log.Println(err)
	}
}
