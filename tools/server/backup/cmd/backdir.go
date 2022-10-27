package main

import (
	"github.com/actliboy/hoper/server/go/lib/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	"log"
	"os"
	"sync"
	"syscall"
	"time"
	"tools/backup"
)

func main() {
	defer initialize.Start(nil, &backup.Dao)()
	backup.Dao.Hoper.Migrator().CreateTable(&backup.File{})
	wg := &sync.WaitGroup{}
	back(wg, "F:/", 0, 0)
	wg.Wait()
}

func back(wg *sync.WaitGroup, dir string, id, level int) {
	entities, err := os.ReadDir(dir)
	if err != nil {
		log.Println(dir, err)
		return
	}
	for _, entity := range entities {

		file := &backup.File{
			Name:  entity.Name(),
			Level: level + 1,
			Pid:   id,
		}

		info, err := entity.Info()
		if err != nil {
			log.Println(dir, err)
		}
		if info != nil {
			wFileSys := info.Sys().(*syscall.Win32FileAttributeData)
			file.CreateTime = time.Unix(0, wFileSys.CreationTime.Nanoseconds())
			file.ModTime = info.ModTime()
		}

		if !entity.IsDir() {
			file.Size = int(info.Size())
			wg.Add(1)
			go func() {
				backup.Dao.Hoper.Create(file)
				wg.Done()
			}()
		} else {
			backup.Dao.Hoper.Create(file)
			back(wg, dir+fs.PathSeparator+entity.Name(), file.Id, file.Level)
		}

	}
}
