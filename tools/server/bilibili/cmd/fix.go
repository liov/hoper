package main

import (
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/initialize"
	"log"
	"os"
	"path"
	"strings"
	"tools/bilibili/config"
	"tools/bilibili/dao"
)

func main() {
	defer initialize.Start(config.Conf, &dao.Dao)()
	fixQuality()
}

func fixRecord() {
	dir := "F:\\B站\\video"
	files, _ := os.ReadDir(dir)
	for _, file := range files {
		cid := strings.Split(file.Name(), "_")[1]
		dao.Dao.Hoper.Table(dao.TableNameVideo).Where(`cid = `+cid).UpdateColumn("record", true)
	}
}

func fixQuality() {
	dir := "D:\\F\\B站\\video"
	files, _ := os.ReadDir(dir)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), "64.flv") || strings.HasSuffix(file.Name(), "80.flv") {
			cid := strings.Split(file.Name(), "_")[1]
			var quality string
			err := dao.Dao.Hoper.Raw(`SELECT data #> '{accept_quality,0}' quality FROM ` + dao.TableNameVideo + ` WHERE cid = ` + cid).Row().Scan(&quality)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Println(file.Name()[len(file.Name())-6 : len(file.Name())-4])
			if quality != file.Name()[len(file.Name())-6:len(file.Name())-4] {
				dao.Dao.Hoper.Table(dao.TableNameVideo).Where(`cid = `+cid).UpdateColumn("record", false)
				os.Remove(path.Join(dir, file.Name()))
			}
		}
	}
}

func remove() {
	dir := "D:\\F\\B站\\video"
	log.Println(path.Dir(dir))
	files, _ := os.ReadDir(dir)
	m := map[string]struct{}{}
	for _, file := range files {
		cid := strings.Split(file.Name(), "_")[1]
		m[cid] = struct{}{}
		/*err := dao.Dao.Hoper.Table(dao.TableNameVideo).Where("cid = "+cid).Update("record", true).Error
		if err != nil {
			log.Println(err)
			return
		}*/
	}
	dir = "F:\\Pictures\\B站"
	files, _ = os.ReadDir(dir)
	for _, file := range files {
		if strings.Contains(file.Name(), "-") {
			cid := strings.Split(file.Name(), "-")[0]
			if _, ok := m[cid]; ok {
				err := os.Remove(path.Join(dir, file.Name()))
				if err != nil {
					log.Println(err)
					return
				}
				log.Println("remove", file.Name())
			}
		}
	}
}
