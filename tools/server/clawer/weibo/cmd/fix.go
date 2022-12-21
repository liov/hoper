package main

import (
	"github.com/liov/hoper/server/go/lib/initialize"
	"github.com/liov/hoper/server/go/lib/utils/fs"
	stringsi "github.com/liov/hoper/server/go/lib/utils/strings"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	claweri "tools/clawer"
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
	zeroTime := time.Time{}
	timer := time.NewTicker(time.Second)
	for _, subdir := range subdirs {
		compsubdir := commondir + fs.PathSeparator + subdir.Name()
		m := make(map[string]time.Time)
		files, _ := os.ReadDir(compsubdir)
		for _, f := range files {
			fname := f.Name()
			info, _ := f.Info()
			strs := strings.Split(fname, "_")

			if len(strs) == 3 {
				date, ok := m[strs[0]+"-"+strs[1]]
				if !ok {
					for i := 0; i < 10; i++ {
						<-timer.C
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
				if date == zeroTime {
					date = info.ModTime()
					m[strs[0]+"-"+strs[1]] = date
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
				newpath := config.Conf.Weibo.DownloadPath + "/" + dir.Path()
				os.MkdirAll(fs.GetDir(newpath), 0666)
				log.Println("rename:", oldpath, newpath)
				err := os.Rename(oldpath, newpath)
				if err != nil {
					log.Println(err)
				}
				dao.Dao.Hoper.Create(dir)
			}
		}
		files, _ = os.ReadDir(compsubdir)
		if len(files) == 0 {
			os.Remove(compsubdir)
		}
	}
}

func rename3() {
	zeroTime := time.Time{}
	commondir := "F:\\Pictures\\pron\\weibo\\debug"
	subdirs, _ := os.ReadDir(commondir)
	timer := time.NewTicker(time.Second)
	for _, subdir := range subdirs {
		compsubdir := commondir + fs.PathSeparator + subdir.Name()
		subdir2s, _ := os.ReadDir(compsubdir)
		for _, subdir2 := range subdir2s {
			compsubdir2 := compsubdir + fs.PathSeparator + subdir2.Name()
			subdir3s, _ := os.ReadDir(compsubdir2)

			for _, subdir3 := range subdir3s {
				compsubdir3 := compsubdir2 + fs.PathSeparator + subdir3.Name()

				m := make(map[string]time.Time)
				files, _ := os.ReadDir(compsubdir3)
				if len(files) == 0 {
					os.Remove(compsubdir3)
				}
				for _, f := range files {
					fname := f.Name()
					info, _ := f.Info()
					strs := strings.Split(fname, "_")

					if len(strs) == 3 {
						date, ok := m[strs[0]+"-"+strs[1]]
						if !ok {
							for i := 0; i < 10; i++ {
								<-timer.C
								weibo, err := rpc.GetLongWeibo(strs[1])
								if err == nil {
									date, _ = time.Parse(time.RubyDate, weibo.CreatedAt)
									m[strs[0]+"-"+strs[1]] = date
									break
								}
								if strings.HasPrefix(err.Error(), "invalid character") {
									date = info.ModTime()
									m[strs[0]+"-"+strs[1]] = date
									break
								}
								log.Println(err)
							}
						}
						if date == zeroTime {
							date = info.ModTime()
							m[strs[0]+"-"+strs[1]] = date
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

						oldpath := compsubdir3 + fs.PathSeparator + fname
						newpath := config.Conf.Weibo.DownloadPath + "/" + dir.Path()
						os.MkdirAll(fs.GetDir(newpath), 0666)
						log.Println("rename:", oldpath, newpath)
						err := os.Rename(oldpath, newpath)
						if err != nil {
							log.Println(err)
						}
						dao.Dao.Hoper.Create(dir)
					}
				}
				files, _ = os.ReadDir(compsubdir3)
				if len(files) == 0 {
					os.Remove(compsubdir)
				}
			}

		}

	}
}

func rename4() {
	commondir := "F:\\Pictures\\pron\\weibo\\debug\\2022"
	subdirs, _ := os.ReadDir(commondir)
	for _, subdir := range subdirs {
		compsubdir := commondir + fs.PathSeparator + subdir.Name()
		subdir2s, _ := os.ReadDir(compsubdir)
		for _, subdir2 := range subdir2s {
			compsubdir2 := compsubdir + fs.PathSeparator + subdir2.Name()

			files, _ := os.ReadDir(compsubdir2)
			for _, f := range files {
				fname := f.Name()
				strs := strings.Split(fname, "_")

				if len(strs) == 4 {
					newDir := "F:\\Pictures\\pron\\weibo\\" + strs[0] + fs.PathSeparator + strs[2][:4]
					os.MkdirAll(newDir, 0666)
					newpath := newDir + fs.PathSeparator + strings.Join([]string{strs[2], strs[0], strs[1], strs[3]}, "_")
					log.Println("rename:", compsubdir2+fs.PathSeparator+fname, newpath)
					err := os.Rename(compsubdir2+fs.PathSeparator+fname, newpath)
					if err != nil {
						log.Println(err)
					}
				}
			}

		}

	}
}

func rename5() {
	commondir := "F:\\Pictures\\pron\\weibo\\2022\\2022-12\\2022-12-20"
	files, _ := os.ReadDir(commondir)
	timer := time.NewTicker(time.Second)
	zeroTime := time.Time{}
	m := make(map[string]time.Time)
	for _, f := range files {
		fname := f.Name()
		info, _ := f.Info()
		strs := strings.Split(fname, "_")
		var dir claweri.Dir
		err := dao.Dao.Hoper.Where(`platform = 4 AND user_id = ` + strs[0] + ` AND key_id_str = '` + strs[1] + `' AND base_url = '` + strs[len(strs)-1] + "'").Find(&dir).Error
		if err != nil {
			log.Println(err)
		}
		if strs[0] == "6537140514" && strs[1] == "4848747278245572" {
			log.Println("到了")
		}
		if dir.Type == 0 || dir.CreatedAt.Sub(dir.PubAt) > time.Hour {
			date, ok := m[strs[0]+"-"+strs[1]]
			if !ok {
				for i := 0; i < 10; i++ {
					<-timer.C
					weibo, err := rpc.GetLongWeibo(strs[1])
					if err == nil {
						date, _ = time.Parse(time.RubyDate, weibo.CreatedAt)
						m[strs[0]+"-"+strs[1]] = date
						if dir.Type == 0 {
							dir.Platform = 4
							dir.UserId, _ = strconv.Atoi(strs[0])
							dir.KeyIdStr = strs[1]
							dir.CreatedAt = info.ModTime()
							dir.PubAt = date
							dir.BaseUrl = strs[len(strs)-1]
							dir.Type = 1
							if strings.HasSuffix(dir.BaseUrl, ".mov") {
								dir.Type = 2
							}
							if strings.HasSuffix(dir.BaseUrl, ".mp4") {
								dir.Type = 3
							}
							dao.Dao.Hoper.Create(&dir)
						} else {
							dao.Dao.Hoper.Table("dir").Where(`platform = 4 AND user_id = `+strs[0]+` AND key_id_str = '`+strs[1]+`' `).Update("pub_at", date)
						}

						break
					} else if strings.HasPrefix(err.Error(), "invalid") {
						if len(strs) == 4 {
							dir.Platform = 4
							dir.UserId, _ = strconv.Atoi(strs[0])
							dir.KeyIdStr = strs[1]
							dir.CreatedAt = info.ModTime()
							dir.PubAt, _ = time.Parse("20060102150405", strs[2])
							dir.BaseUrl = strs[len(strs)-1]
							dir.Type = 1
							if strings.HasSuffix(dir.BaseUrl, ".mov") {
								dir.Type = 2
							}
							if strings.HasSuffix(dir.BaseUrl, ".mp4") {
								dir.Type = 3
							}
							dao.Dao.Hoper.Create(&dir)
							break
						}
					}
				}
			}
			if date == zeroTime {
				date = dir.CreatedAt
				m[strs[0]+"-"+strs[1]] = date
				dao.Dao.Hoper.Table("dir").Where(`platform = 4 AND user_id = `+strs[0]+` AND key_id_str = '`+strs[1]+"'").Update("pub_at", date)
			}
			dir.PubAt = date

		}
		log.Println("rename:", commondir+fs.PathSeparator+fname, "F:\\Pictures\\pron\\weibo\\"+dir.Path())
		err = os.Rename(commondir+fs.PathSeparator+fname, "F:\\Pictures\\pron\\weibo\\"+dir.Path())
		if err != nil {
			log.Println(err)
		}
	}
}
