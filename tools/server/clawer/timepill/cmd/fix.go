package main

import (
	"github.com/liov/hoper/server/go/lib/initialize"
	"github.com/liov/hoper/server/go/lib/utils/fs"
	timei "github.com/liov/hoper/server/go/lib/utils/time"
	"log"
	"os"
	"strconv"
	claweri "tools/clawer"
	"tools/clawer/timepill"
	"tools/clawer/timepill/model"
)

func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()
	rename2()
}

func rename2() {
	commondir := "D:\\F\\timepill\\2022\\2022-12"
	subdirs, _ := os.ReadDir(commondir)

	for _, subdir := range subdirs {
		compsubdir := commondir + fs.PathSeparator + subdir.Name()
		subdirs12, _ := os.ReadDir(compsubdir)
		for _, subdir1 := range subdirs12 {
			compsubdir2 := compsubdir + fs.PathSeparator + subdir1.Name()
			files, _ := os.ReadDir(compsubdir2)
			for _, f := range files {
				fname := f.Name()
				info, _ := f.Info()

				userId, _ := strconv.Atoi(subdir.Name())
				id, _ := strconv.Atoi(subdir1.Name())

				baseUrl := fname
				var diary model.Diary
				err := timepill.Dao.Hoper.Where(`user_id = ? AND id = ?`, userId, id).First(&diary).Error
				if err != nil {
					log.Println(err)
				}
				pubAt, _ := timei.Parse(timei.TimeFormatPostgresDB, diary.Created)

				oldpath := compsubdir2 + fs.PathSeparator + fname
				num := userId / 10000
				newdir := "D:\\F\\timepill\\debug\\" + strconv.Itoa(num) + "-" + strconv.Itoa(num+1)
				os.MkdirAll(newdir+fs.PathSeparator+subdir.Name()+fs.PathSeparator+diary.Created[:4], 0666)

				dir := &claweri.Dir{
					Platform:  2,
					UserId:    userId,
					KeyId:     diary.Id,
					BaseUrl:   baseUrl,
					Type:      1,
					PubAt:     pubAt,
					CreatedAt: info.ModTime(),
				}
				if err != nil {
					dir.KeyIdStr = "unknown"
				}
				log.Println("rename:", oldpath, newdir+fs.PathSeparator+dir.Path())
				err = os.Rename(oldpath, newdir+fs.PathSeparator+dir.Path())
				if err != nil {
					log.Println(err)
				}

				//timepill.Dao.Hoper.Create(dir)
			}
			files, _ = os.ReadDir(compsubdir2)
			if len(files) == 0 {
				err := os.Remove(compsubdir2)
				if err != nil {
					log.Println(err)
				}
			}
		}
		subdirs12, _ = os.ReadDir(compsubdir)
		if len(subdirs12) == 0 {
			err := os.Remove(compsubdir)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
