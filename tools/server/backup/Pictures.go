package main

import (
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	"github.com/jlaffaye/ftp"
	"log"
	"sort"
)

func Pietures(c *ftp.ServerConn) {
	err := Copy(c, "/Pictures/Twitter", BackUpDisk+"pron\\Twitter")
	if err != nil {
		log.Println(err)
	}

	err = Copy(c, "/Pictures/douyin", BackUpDisk+"douyin")
	if err != nil {
		log.Println(err)
	}

	err = Copy(c, "/Movies/TwDown", BackUpDiskPron+"TwDown")
	if err != nil {
		log.Println(err)
	}
	err = Copy(c, "/Movies/Telegram", BackUpDiskPron+"Telegram")
	if err != nil {
		log.Println(err)
	}
	err = Copy(c, "/Pictures/CoolMarket", BackUpDisk+"CoolMarket")
	if err != nil {
		log.Println(err)
	}
	err = Copy(c, "/Pictures/weibo", BackUpDiskPron+"weibo")
	if err != nil {
		log.Println(err)
	}
	err = Copy(c, "/Pictures/微博动图", BackUpDiskPron+"微博动图")
	if err != nil {
		log.Println(err)
	}
	err = Copy(c, "/Pictures/Vipaccount", BackUpDisk+"Vipaccount")
	if err != nil {
		log.Println(err)
	}
	err = Copy(c, "/nico/downloadImages", BackUpDisk+"nico\\downloadImages")
	if err != nil {
		log.Println(err)
	}
	err = Copy(c, "/soul/Pictures", BackUpDisk+"soul\\Pictures")
	if err != nil {
		log.Println(err)
	}
	err = Copy(c, "/Download/IDMP/Videos", BackUpDiskPron+"Videos")
	if err != nil {
		log.Println(err)
	}
	err = Copy(c, "/DCIM/weibo", BackUpDiskPron+"Share")
	if err != nil {
		log.Println(err)
	}
	err = Copy(c, "/DCIM/weibovideo/", BackUpDiskPron+"Share")
	if err != nil {
		log.Println(err)
	}
}

func Copy(c *ftp.ServerConn, src, dst string) error {
	lastFile, m, err := fs.LastFile(dst)
	if err != nil {
		return err
	}
	log.Println(lastFile.Name())

	list, err := c.List(src)
	if err != nil {
		return err
	}
	log.Println(len(list))
	sort.Sort(Entities(list))
	var lastIdx int
	for i, item := range list {
		if item.Name == lastFile.Name() {
			lastIdx = i
			break
		}
	}
	if lastIdx == 0 {
		for i := 0; i < len(list); i++ {
			item := list[i]

			if _, ok := m[item.Name]; ok {
				lastIdx = i
				break
			}
		}
	}
	if lastIdx == 0 {
		lastIdx = len(list)
	}
	for i := lastIdx - 1; i >= 0; i-- {
		item := list[i]
		if item.Type != ftp.EntryTypeFile {
			continue
		}
		resp, err := c.Retr(src + "/" + item.Name)
		if err != nil {
			return err
		}

		err = fs.Copy(dst+sep+item.Name, resp)
		if err != nil {
			return err
		}
		log.Println("copy file: ", dst+sep+item.Name)
		resp.Close()
	}
	return nil
}

type Rule func(string) string
