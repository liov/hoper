package main

import (
	"github.com/liov/hoper/server/go/lib/initialize"
	"github.com/liov/hoper/server/go/lib/utils/fs"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	claweri "tools/clawer"
	"tools/clawer/timepill"
	"tools/clawer/weibo/rpc"
)

func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()
}

func rename2() {
	commondir := "D:\\F\\timepill\\2010_"
	subdirs, _ := os.ReadDir(commondir)
	zeroTime := time.Time{}
	for _, subdir := range subdirs {
		compsubdir := commondir + fs.PathSeparator + subdir.Name()
		m := make(map[string]time.Time)
		files, _ := os.ReadDir(compsubdir)
		for _, f := range files {
			fname := f.Name()
			info, _ := f.Info()
			strs := strings.Split(fname, "_")

			for i := 0; i < 10; i++ {
				weibo, err := rpc.GetLongWeibo(strs[1])
				if err == nil {
					date, _ = time.Parse(time.RubyDate, weibo.CreatedAt)

					m[strs[0]+"-"+strs[1]] = date
					break
				}
				if strings.HasPrefix(err.Error(), "json.Unmarshal error:invalid character") {
					date = info.ModTime()
					m[strs[0]+"-"+strs[1]] = date
					break
				}
				log.Println(err)
			}
		}

		userId, _ := strconv.Atoi(strs[0])
		dir := &claweri.Dir{
			Platform:  4,
			UserId:    userId,
			KeyIdStr:  strs[1],
			BaseUrl:   strs[2],
			Type:      1,
			PubAt:     date,
			CreatedAt: info.ModTime(),
		}
		if strings.HasSuffix(strs[2], ".mov") {
			dir.Type = 2
		}
		oldpath := compsubdir + fs.PathSeparator + fname
		newpath := timepill.Conf.TimePill.PhotoPath + "/" + dir.Path()
		os.MkdirAll(fs.GetDir(newpath), 0666)
		log.Println("rename:", oldpath, newpath)
		err := os.Rename(oldpath, newpath)
		if err != nil {
			log.Println(err)
		}
		timepill.Dao.Hoper.Create(dir)
	}
	files, _ = os.ReadDir(compsubdir)
	if len(files) == 0 {
		os.Remove(compsubdir)
	}
}
}
