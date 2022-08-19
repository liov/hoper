package main

import (
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	"github.com/jlaffaye/ftp"
	"log"
	"os"
	"sort"
	"strings"
)

const BackUpDisk = "F:\\Pictures\\"
const BackUpDiskPron = "F:\\Pictures\\pron\\"

const sep = string(os.PathSeparator)

func DCIM(c *ftp.ServerConn) {
	jpdir := BackUpDiskPron + "pic\\jiepai"
	jplastFile, err := fs.LastFile(jpdir)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(jplastFile.Name())

	xhsdir := BackUpDisk + "XHS"
	xhslastFile, err := fs.LastFile(xhsdir)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(xhslastFile.Name())

	dydir := BackUpDisk + "douyin"
	dylastFile, err := fs.LastFile(dydir)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(dylastFile.Name())

	dyvdir := BackUpDisk + "douyin_video"
	dyvlastFile, err := fs.LastFile(dyvdir)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(dyvlastFile.Name())

	lastFile := jplastFile
	for _, file := range []os.FileInfo{xhslastFile, dylastFile, dyvlastFile} {
		if file.ModTime().After(lastFile.ModTime()) {
			lastFile = xhslastFile
		}
	}

	srcdir := "/DCIM/Camera"
	list, err := c.List(srcdir)
	if err != nil {
		log.Println(err)
		return
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
	for i := lastIdx; i > 0; i-- {
		item := list[i]
		if item.Type != ftp.EntryTypeFile {
			continue
		}
		resp, err := c.Retr(srcdir + "/" + item.Name)
		if err != nil {
			log.Println(err)
			return
		}
		//log.Println(item)
		switch {
		case strings.HasPrefix(item.Name, "IMG"), strings.HasPrefix(item.Name, "MVIMG"), strings.HasPrefix(item.Name, "VID"):
			err = fs.Copy(jpdir+sep+item.Name, resp)
		case strings.HasPrefix(item.Name, "XHS"):
			err = fs.Copy(xhsdir+sep+item.Name, resp)
		default:
			if strings.HasSuffix(item.Name, "mp4") {
				err = fs.Copy(dyvdir+sep+item.Name, resp)
			} else {
				err = fs.Copy(dydir+sep+item.Name, resp)
			}
		}
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("copy file: ", item.Name)
		resp.Close()
	}
}
