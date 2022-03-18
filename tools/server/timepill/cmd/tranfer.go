package main

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	surl "net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"tools/timepill"
)

func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()
	induction()
}

func getDir(userId int, url, created string) (string, string) {
	if url == "" {
		return "", ""
	}
	URL, _ := surl.Parse(url)
	suffixpath := URL.Path
	if strings.HasSuffix(URL.Path, "!large") {
		suffixpath = strings.TrimSuffix(URL.Path, "!large")
	} else if strings.Contains(URL.Path, "photos") {
		suffixpath = strings.Split(URL.Path, "photos")[1]
	}
	date := suffixpath[1:11]
	if strings.Contains(date, "/") {
		date = created[0:10]
	}
	oldpath := timepill.Conf.TimePill.PhotoPath + "/" + suffixpath[1:5] + suffixpath
	newpath := timepill.Conf.TimePill.PhotoPath + "/" + strconv.Itoa(userId) + "/" + date + "-" + path.Base(suffixpath)
	return oldpath, newpath
}

func getFileName(url string) string {
	if url == "" {
		return ""
	}
	URL, _ := surl.Parse(url)
	suffixpath := URL.Path
	if strings.HasSuffix(URL.Path, "!large") {
		suffixpath = strings.TrimSuffix(URL.Path, "!large")
	}
	return path.Base(suffixpath)
}

func transfer(diary *timepill.TinyDiary) {
	oldpath, newpath := getDir(diary.UserId, diary.PhotoUrl, diary.Created)
	_, err := os.Stat(oldpath)
	if os.IsNotExist(err) {
		timepill.DownloadPic(diary.UserId, diary.PhotoUrl, diary.Created)
		return
	}
	_, err = os.Stat(newpath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path.Dir(newpath), 0666)
		if err != nil {
			log.Error(err)
		}
	}
	err = os.Rename(oldpath, newpath)
	if err != nil {
		log.Error(err)
	}
}

func induction() {
	dir := timepill.Conf.TimePill.PhotoPath + "/"
	dirInfos, err := os.ReadDir(dir)
	if err != nil {
		log.Error(err)
	}
	for i := range dirInfos {
		if dirInfos[i].IsDir() {
			if strings.Contains(dirInfos[i].Name(), "-") {
				continue
			}
			num, _ := strconv.Atoi(dirInfos[i].Name())
			num /= 10000
			outdir := dir + strconv.Itoa(num) + "-" + strconv.Itoa(num+1) + "/" + dirInfos[i].Name() + "/"
			_, err = os.Stat(outdir)
			if os.IsNotExist(err) {
				os.MkdirAll(outdir, 0666)
			}
			subDir := dir + dirInfos[i].Name() + "/"
			fileInfos, err := os.ReadDir(subDir)
			if err != nil {
				log.Error(err)
			}
			for j := range fileInfos {
				if !fileInfos[j].IsDir() {
					err = os.Rename(subDir+fileInfos[j].Name(), outdir+fileInfos[j].Name())
					if err != nil {
						log.Error(err)
					}
				}
			}
			fileInfos, err = os.ReadDir(subDir)
			if err != nil {
				log.Error(err)
			}
			if len(fileInfos) == 0 {
				err = os.Remove(subDir)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
}
