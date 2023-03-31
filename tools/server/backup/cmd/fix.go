package main

import (
	"github.com/hopeio/pandora/initialize"
	"github.com/hopeio/pandora/utils/fs"
	"github.com/hopeio/pandora/utils/log"
	"os"
	"strings"
	"syscall"
	"time"
	"tools/backup"
)

func main() {
	defer initialize.Start(nil, &backup.Dao)()
	//aliyun()
	type Dir struct {
		D1, D2, D3 string
	}

	//normal(backup.BackUpDiskPron + "OracleBuBu")
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

func aliyun() {
	aliyun := "F:\\Pictures\\pron\\pic\\aliyun"
	dir := "F:\\Pictures\\pron\\pic\\WeiXin"
	fileMap := make(map[string]struct{})
	getFileMap(dir, fileMap)
	entities, _ := os.ReadDir(aliyun)
	if len(entities) == 0 {
		return
	}
	for _, entity := range entities {
		if _, ok := fileMap[entity.Name()]; ok {
			os.Remove(aliyun + fs.PathSeparator + entity.Name())
		} else {
			var file backup.File
			err := backup.Dao.Hoper.Where("name = ?", entity.Name()).First(&file).Error
			if err != nil {
				log.Error(err)
				continue
			}
			oldpath := aliyun + fs.PathSeparator + entity.Name()
			date := file.ModTime.Format("200601")
			newDir := dir + fs.PathSeparator + date
			_, err = os.Stat(newDir)
			if os.IsNotExist(err) {
				err = os.Mkdir(newDir, 0666)
				if err != nil {
					log.Error(err)
				}
			}
			newPath := newDir + fs.PathSeparator + entity.Name()
			_, err = os.Stat(newPath)
			if os.IsNotExist(err) {
				log.Info("rename:", newPath)
				err = os.Rename(oldpath, newPath)
				if err != nil {
					log.Error(err)
				}
			} else {
				log.Info("delete:", oldpath)
				os.Remove(oldpath)
			}
		}
	}

}

func getFileMap(dir string, fileMap map[string]struct{}) {
	entities, _ := os.ReadDir(dir)
	for _, entity := range entities {
		if entity.IsDir() {
			getFileMap(dir+fs.PathSeparator+entity.Name(), fileMap)
		}
		fileMap[entity.Name()] = struct{}{}
	}
}

const sql = `SELECT b3.NAME  d1,b2.NAME d2,b1.NAME d3 FROM "public"."bak_dir1" b1 LEFT JOIN "public"."bak_dir1" b2 ON b1.pid = b2.ID LEFT JOIN "public"."bak_dir1" b3 ON b2.pid = b3.ID WHERE b1."name" = ?`
