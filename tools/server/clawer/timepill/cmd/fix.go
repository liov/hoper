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
	rename2()
}

func rename() {
	commondir := "D:\\F\\timepill\\debug\\2022\\2022-12\\2022-12-21"
	files, _ := os.ReadDir(commondir)
	for _, f := range files {
		strs := strings.Split(f.Name(), "_")
		userId, _ := strconv.Atoi(strs[0])
		num := userId / 10000
		newdir := "D:\\F\\timepill\\debug\\" + strconv.Itoa(num) + "-" + strconv.Itoa(num+1) + fs.PathSeparator + strs[0] + fs.PathSeparator + "2022"
		os.MkdirAll(newdir, 0666)
		log.Println("rename:", commondir+fs.PathSeparator+f.Name(), newdir+fs.PathSeparator+strings.Join([]string{strs[2], strs[0], strs[1], strs[3]}, "_"))
		os.Rename(commondir+fs.PathSeparator+f.Name(), newdir+fs.PathSeparator+strings.Join([]string{strs[2], strs[0], strs[1], strs[3]}, "_"))
	}

}

func rename2() {
	commondir := "D:\\F\\timepill\\2021_"
	subdirs, _ := os.ReadDir(commondir)
	zeroTime := time.Time{}
	for _, subdir := range subdirs {
		compsubdir := commondir + fs.PathSeparator + subdir.Name()
		files, _ := os.ReadDir(compsubdir)
		for _, f := range files {
			fname := f.Name()
			info, _ := f.Info()
			strs := strings.Split(fname, "_")

			userId, _ := strconv.Atoi(strs[0])
			date, _ := timei.Parse(timei.DateFormat, subdir.Name())

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
				continue
			}
			pubAt, _ := timei.Parse(timei.TimeFormatPostgresDB, diary.Created)
			if pubAt == zeroTime {
				pubAt, _ = timei.Parse(timei.TimeFormatPostgresDB, diary.Updated)
			}
			if pubAt == zeroTime {
				pubAt = date
			}
			oldpath := compsubdir + fs.PathSeparator + fname
			num := userId / 10000
			newdir := "D:\\F\\timepill\\debug\\" + strconv.Itoa(num) + "-" + strconv.Itoa(num+1)
			os.MkdirAll(newdir+fs.PathSeparator+strs[0]+fs.PathSeparator+diary.Created[:4], 0666)

			dir := &claweri.Dir{
				Platform:  2,
				UserId:    userId,
				KeyId:     diary.Id,
				BaseUrl:   baseUrl,
				Type:      1,
				PubAt:     pubAt,
				CreatedAt: info.ModTime(),
			}

			log.Println("rename:", oldpath, newdir+fs.PathSeparator+dir.Path())
			err = os.Rename(oldpath, newdir+fs.PathSeparator+dir.Path())
			if err != nil {
				log.Println(err)
			}

			timepill.Dao.Hoper.Create(dir)
		}
		files, _ = os.ReadDir(compsubdir)
		if len(files) == 0 {
			err := os.Remove(compsubdir)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
