package main

import (
	"github.com/liov/hoper/server/go/lib/initialize"
	"github.com/liov/hoper/server/go/lib/utils/fs"
	timei "github.com/liov/hoper/server/go/lib/utils/time"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	claweri "tools/clawer"
	"tools/clawer/timepill"
	"tools/clawer/timepill/model"
)

func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()
	rename3()
}

func rename3() {
	commondir := "D:\\F\\timepill"
	subdirs, _ := os.ReadDir(commondir)
	zeroTime := time.Time{}
	for _, subdir := range subdirs {
		if strings.Contains(subdir.Name(), "-") {
			compsubdir := commondir + fs.PathSeparator + subdir.Name()
			subdirs2, _ := os.ReadDir(compsubdir)
			for _, subdir2 := range subdirs2 {
				compsubdir2 := compsubdir + fs.PathSeparator + subdir2.Name()
				files, _ := os.ReadDir(compsubdir2)
				for _, f := range files {
					fname := f.Name()
					info, _ := f.Info()
					strs := strings.Split(fname, "_")

					userId, _ := strconv.Atoi(subdir2.Name())
					date, _ := timei.Parse(timei.DateFormat, strs[0])

					var baseUrl string
					if len(strs) == 3 {
						baseUrl = strs[2]
					} else if len(strs) == 2 {
						baseUrl = strs[1]
					}
					var diary model.Diary
					err := timepill.Dao.Hoper.Where(`user_id = ? AND created BETWEEN ? AND ? AND photo_url LIKE ?`, userId, date, date.AddDate(0, 0, 1), "%"+baseUrl+"%").First(&diary).Error
					if err != nil {
						log.Println(err)
						diary.Created = strs[0]
					}
					pubAt, _ := timei.Parse(timei.TimeFormatPostgresDB, diary.Created)
					if pubAt == zeroTime {
						pubAt, _ = timei.Parse(timei.TimeFormatPostgresDB, diary.Updated)
					}
					if pubAt == zeroTime {
						pubAt = date
					}
					oldpath := compsubdir2 + fs.PathSeparator + fname
					num := userId / 10000
					newdir := "D:\\F\\timepill\\debug\\" + strconv.Itoa(num) + "-" + strconv.Itoa(num+1)

					dir := &claweri.Dir{
						Platform:  2,
						UserId:    userId,
						KeyId:     diary.Id,
						BaseUrl:   baseUrl,
						Type:      1,
						PubAt:     pubAt,
						CreatedAt: info.ModTime(),
					}
					if fs.NotExist(newdir + fs.PathSeparator + dir.Path()) {
						os.MkdirAll(newdir+fs.PathSeparator+subdir2.Name()+fs.PathSeparator+diary.Created[:4], 0666)
						log.Println("rename:", oldpath, newdir+fs.PathSeparator+dir.Path())
						err = os.Rename(oldpath, newdir+fs.PathSeparator+dir.Path())
						if err != nil {
							log.Println(err)
						}
					} else {
						os.Remove(oldpath)
					}

					timepill.Dao.Hoper.Create(dir)
				}
				files, _ = os.ReadDir(compsubdir2)
				if len(files) == 0 {
					err := os.Remove(compsubdir2)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}

	}
}
