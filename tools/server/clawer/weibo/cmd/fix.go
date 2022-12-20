package main

import (
	"github.com/liov/hoper/server/go/lib/initialize"
	"github.com/liov/hoper/server/go/lib/utils/fs"
	stringsi "github.com/liov/hoper/server/go/lib/utils/strings"
	timei "github.com/liov/hoper/server/go/lib/utils/time"
	"log"
	"os"
	"strings"
	"time"
	"tools/clawer/weibo/config"
	"tools/clawer/weibo/dao"
	"tools/clawer/weibo/rpc"
)

func main() {
	defer initialize.Start(config.Conf, &dao.Dao)()
	rename2()
}

func rename() {
	commondir := "F:\\Pictures\\pron\\weibo\\pic"
	fs.RangeDir(commondir, func(subDir string, entry os.DirEntry) error {
		fileName := entry.Name()
		if strings.HasSuffix(fileName, "mov") {
			parts := strings.Split(fileName, "_")
			for _, part := range parts {
				if strings.HasSuffix(part, "mov") {
					parts[2] = stringsi.CountdownCutoff(part, "%2F")
					break
				}
			}
			parts[2] = stringsi.CountdownCutoff(parts[2], "%2F")
			path := subDir + fs.PathSeparator + fileName
			newPath := subDir + fs.PathSeparator + strings.Join(parts[:3], "_")
			log.Println("rename:", path, newPath)
			err := os.Rename(path, newPath)
			if err != nil {
				log.Println(err)
			}
		}
		return nil
	})
}

func rename2() {
	commondir := "F:\\Pictures\\pron\\weibo\\pic"
	subdirs, _ := os.ReadDir(commondir)
	timer := time.NewTicker(time.Second)
	for _, subdir := range subdirs {
		compsubdir := commondir + fs.PathSeparator + subdir.Name()
		m := make(map[string]string)
		files, _ := os.ReadDir(compsubdir)
		for _, f := range files {
			fname := f.Name()
			strs := strings.Split(fname, "_")

			if len(strs) == 3 {
				date, ok := m[strs[0]+"-"+strs[1]]
				if !ok {
					for i := 0; i < 10; i++ {
						<-timer.C
						weibo, err := rpc.GetLongWeibo(strs[1])
						if err == nil {
							createdAt, _ := time.Parse(time.RubyDate, weibo.CreatedAt)
							date = createdAt.Format(timei.DateFormat)
							m[strs[0]+"-"+strs[1]] = date
							break
						}
						if strings.HasPrefix(err.Error(), "invalid character") {
							info, _ := f.Info()
							date = info.ModTime().Format(timei.DateFormat)
							m[strs[0]+"-"+strs[1]] = date
							break
						}
						log.Println(err)
					}
				}
				if date == "" {
					info, _ := f.Info()
					date = info.ModTime().Format(timei.DateFormat)
					m[strs[0]+"-"+strs[1]] = date
				}

				oldpath := compsubdir + fs.PathSeparator + fname
				newdir := strings.Join([]string{config.Conf.Weibo.DownloadPath, date[:4], date[:7], date}, "/")
				os.MkdirAll(newdir, 0666)
				newpath := newdir + "/" + fname
				log.Println("rename:", oldpath, newpath)
				err := os.Rename(oldpath, newpath)
				if err != nil {
					log.Println(err)
				}
			}
		}
		files, _ = os.ReadDir(compsubdir)
		if len(files) == 0 {
			os.Remove(compsubdir)
		}
	}
}
