package backup

import (
	"github.com/jlaffaye/ftp"
	"github.com/liov/hoper/server/go/lib/utils/fs"
	"log"
	"sort"
)

func Pietures(c *ftp.ServerConn) {
	err := Copy(c, "/Pictures/Twitter", BackUpDisk+"pron\\Twitter", true)
	if err != nil {
		log.Println(err)
	}

	err = Copy(c, "/Pictures/douyin", BackUpDisk+"douyin", true)
	if err != nil {
		log.Println(err)
	}

	err = Copy(c, "/Movies/TwDown", BackUpDiskPron+"TwDown", true)
	if err != nil {
		log.Println(err)
	}
	err = Copy(c, "/Movies/Telegram", BackUpDiskPron+"Telegram", false)
	if err != nil {
		log.Println(err)
	}
	err = Copy(c, "/Pictures/CoolMarket", BackUpDisk+"CoolMarket", false)
	if err != nil {
		log.Println(err)
	}
	err = Copy(c, "/Pictures/weibo", BackUpDiskPron+"weibo/杂集", true)
	if err != nil {
		log.Println(err)
	}
	err = Copy(c, "/Pictures/微博动图", BackUpDiskPron+"weibo/微博动图", true)
	if err != nil {
		log.Println(err)
	}
	err = Copy(c, "/Pictures/Vipaccount", BackUpDisk+"Vipaccount", false)
	if err != nil {
		log.Println(err)
	}
	err = Copy(c, "/nico/downloadImages", BackUpDisk+"nico\\downloadImages", false)
	if err != nil {
		log.Println(err)
	}
	err = Copy(c, "/soul/Pictures", BackUpDisk+"soul\\Pictures", false)
	if err != nil {
		log.Println(err)
	}
	err = Copy(c, "/Download/IDMP/Videos", BackUpDiskPron+"Videos", false)
	if err != nil {
		log.Println(err)
	}
	err = Copy(c, "/DCIM/weibo", BackUpDiskPron+"Share", false)
	if err != nil {
		log.Println(err)
	}
	err = Copy(c, "/DCIM/weibovideo", BackUpDiskPron+"Share", false)
	if err != nil {
		log.Println(err)
	}
	err = Copy(c, "/DCIM/Screenshots", BackUpDiskPron+"pic\\weibo", false)
	if err != nil {
		log.Println(err)
	}
}

func Copy(c *ftp.ServerConn, src, dst string, date bool) error {
	lastFileFun := fs.LastFile
	if date {
		lastFileFun = LastFile
	}
	lastFile, m, err := lastFileFun(dst)
	if err != nil {
		return err
	}
	log.Println(dst, lastFile.Name())

	list, err := c.List(src)
	if err != nil {
		return err
	}
	log.Println(len(list))
	sort.Sort(Entities(list))
	var lastIdx int
	var findLast bool
	for i, item := range list {
		if item.Name == lastFile.Name() {
			lastIdx = i
			findLast = true
			break
		}
	}
	if lastIdx == 0 && !findLast {
		for i := 0; i < len(list); i++ {
			item := list[i]

			if _, ok := m[item.Name]; ok {
				lastIdx = i
				findLast = true
				break
			}
		}
	}
	if lastIdx == 0 && !findLast {
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
		filename := dst + sep + item.Name
		if date {
			filename = dst + sep + item.Time.Format("200601") + sep + item.Name
		}
		err = fs.Copy(filename, resp)
		if err != nil {
			return err
		}
		log.Println("copy file: ", filename)
		resp.Close()
	}
	return nil
}

type Rule func(string) string
