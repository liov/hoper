package main

import (
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"os"
	"strings"
	"syscall"
	"time"
	"tools/backup"
)

func main() {
	normal(backup.BackUpDiskPron + "OracleBuBu")
}

func dcim() {
	opt := backup.BackUpDiskPron + "pic\\jiepai"
	entities, err := os.ReadDir(opt)
	if len(entities) == 0 {
		log.Error(err)
	}
	mkdir := mkdir{}
	for _, entity := range entities {
		if entity.IsDir() {
			continue
		}
		if strings.HasPrefix(entity.Name(), "IMG_") {
			//info,_:=entity.Info()
			mkdir.Rename(opt, entity.Name(), entity.Name()[4:10])
		} else if strings.HasPrefix(entity.Name(), "IMG") {
			mkdir.Rename(opt, entity.Name(), entity.Name()[3:9])
		}
		if strings.HasPrefix(entity.Name(), "VID_") {
			//info,_:=entity.Info()
			mkdir.Rename(opt, entity.Name(), entity.Name()[4:10])
		}
		if strings.HasPrefix(entity.Name(), "MVIMG_") {
			//info,_:=entity.Info()
			mkdir.Rename(opt, entity.Name(), entity.Name()[6:12])
		}
	}
}

func normal(opt string) {
	entities, err := os.ReadDir(opt)
	if len(entities) == 0 {
		log.Error(err)
	}
	mkdir := mkdir{}
	for _, entity := range entities {
		if entity.IsDir() {
			continue
		}

		info, _ := entity.Info()
		mkdir.Rename(opt, entity.Name(), info.ModTime().Format("200601"))

	}
}

type mkdir struct {
	date [30][13]bool
}

// 200601
func (d *mkdir) getDate(date string) (year, month int) {
	return int((date[2]-'0')*10 + (date[3] - '0')), int((date[4]-'0')*10 + (date[5] - '0'))
}

func (d *mkdir) Rename(opt, filename, date string) {
	year, month := d.getDate(date)
	dir := opt + "\\" + date
	if !d.date[year][month] {
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err = os.Mkdir(dir, 0666)
			if err != nil {
				log.Error(err)
			}
		}
		d.date[year][month] = true
	}
	log.Info("rename:", dir+"\\"+filename)
	err := os.Rename(opt+"\\"+filename, dir+"\\"+filename)
	if err != nil {
		log.Error(err)
	}
}

func fix2(dir string) {
	entities, _ := os.ReadDir(dir)
	for _, entity := range entities {
		if entity.IsDir() {
			continue
		}
		info, _ := entity.Info()
		log.Info(info)
		winFile := info.Sys().(*syscall.Win32FileAttributeData)
		log.Info(time.Unix(0, winFile.CreationTime.Nanoseconds()), time.Unix(0, winFile.LastWriteTime.Nanoseconds()))

	}
}
