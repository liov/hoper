package main

import (
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	//defer initialize.Start(config.Conf, &dao.Dao)()
	dir := "G:\\B站\\video"
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
			cid := strings.Split(file.Name(), "_")[0]
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
