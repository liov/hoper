package main

import (
	"github.com/actliboy/hoper/server/go/lib/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	surl "net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"tools/timepill"
	"tools/timepill/model"
)

func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()
	tx()
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
	newpath := timepill.Conf.TimePill.PhotoPath + "/" + strconv.Itoa(userId) + "/" + date + "_" + path.Base(suffixpath)
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
	} else if strings.Contains(URL.Path, "photos") {
		suffixpath = strings.Split(URL.Path, "photos")[1]
	}
	return path.Base(suffixpath)
}

func transfer(year string) {
	dir := timepill.Conf.TimePill.PhotoPath + "/" + year
	dirInfos, err := os.ReadDir(dir)
	if err != nil {
		log.Error(err)
	}
	for i := range dirInfos {
		var diaries []*model.TinyDiary
		date := dirInfos[i].Name()
		timepill.Dao.Hoper.Table(model.DiaryTableName).Select(`user_id,photo_url,created`).Where(`type = 2 AND created BETWEEN ? AND ? `, date+" 00:00:00", date+" 23:59:59").Scan(&diaries)
		for _, diary := range diaries {
			oldpath, newpath := getDir(diary.UserId, diary.PhotoUrl, diary.Created)
			_, err := os.Stat(oldpath)
			if os.IsNotExist(err) {
				err = timepill.DownloadPic(diary.UserId, diary.PhotoUrl, diary.Created)
				if err != nil {
					log.Error(err)
				}
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

func tx() {
	dir := timepill.Conf.TimePill.PhotoPath + "/"
	dirInfos, err := os.ReadDir(dir)
	if err != nil {
		log.Error(err)
	}
	for i := range dirInfos {
		if strings.Contains(dirInfos[i].Name(), "-") || strings.Contains(dirInfos[i].Name(), "_") {
			continue
		}
		subDir := dir + dirInfos[i].Name() + "/"
		fileInfos, err := os.ReadDir(subDir)
		if err != nil {
			log.Error(err)
		}
		for j := range fileInfos {
			if !fileInfos[j].IsDir() {
				filename := fileInfos[j].Name()
				year := filename[0:4] + "_"
				_, err := os.Stat(dir + year)
				if os.IsNotExist(err) {
					os.Mkdir(dir+year, 0666)
				}
				date := filename[0:10]
				_, err = os.Stat(dir + year + "/" + date)
				if os.IsNotExist(err) {
					os.Mkdir(dir+year+"/"+date, 0666)
				}
				err = timepill.CopyDatePic(subDir+filename, date, dirInfos[i].Name(), filename[11:])
				if err != nil {
					log.Error(err)
				}
				num, _ := strconv.Atoi(dirInfos[i].Name())
				num /= 10000
				outdir := dir + strconv.Itoa(num) + "-" + strconv.Itoa(num+1) + "/" + dirInfos[i].Name() + "/"
				_, err = os.Stat(outdir)
				if os.IsNotExist(err) {
					os.MkdirAll(outdir, 0666)
				}
				err = os.Rename(subDir+filename, outdir+filename[0:10]+"_"+filename[11:])
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
