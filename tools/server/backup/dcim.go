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
	jplastFile, jpm, err := fs.LastFile(jpdir)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(jplastFile.Name())

	xhsdir := BackUpDisk + "XHS"
	xhslastFile, xhsm, err := fs.LastFile(xhsdir)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(xhslastFile.Name())

	dydir := BackUpDisk + "douyin"
	dylastFile, dym, err := fs.LastFile(dydir)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(dylastFile.Name())

	dyvdir := BackUpDisk + "douyin_video"
	dyvlastFile, dyvm, err := fs.LastFile(dyvdir)
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
	if lastIdx == 0 {
		for i := 0; i < len(list); i++ {
			item := list[i]
			if _, ok := jpm[item.Name]; ok {
				lastIdx = i
				break
			}
			if _, ok := xhsm[item.Name]; ok {
				lastIdx = i
				break
			}
			if _, ok := dym[item.Name]; ok {
				lastIdx = i
				break
			}
			if _, ok := dyvm[item.Name]; ok {
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
		resp, err := c.Retr(srcdir + "/" + item.Name)
		if err != nil {
			log.Println(err)
			return
		}
		dst := jpdir
		switch {
		case strings.HasPrefix(item.Name, "IMG"), strings.HasPrefix(item.Name, "MVIMG"), strings.HasPrefix(item.Name, "VID"):
			dst = jpdir
		case strings.HasPrefix(item.Name, "XHS"):
			dst = xhsdir
		default:
			if strings.HasSuffix(item.Name, "mp4") {
				dst = dyvdir
			} else {
				dst = dydir
			}
		}
		err = fs.Copy(dst+sep+item.Name, resp)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("copy file: ", dst+sep+item.Name)
		resp.Close()
	}
}
