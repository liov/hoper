package main

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client"
	surl "net/url"
	"path"
	"strings"
	"tools/timepill"
)

func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()
	date := "2022-02-16"
	var diaries []*timepill.TinyDiary
	timepill.Dao.Hoper.Table(`diary`).Select(`user_id,photo_url,created`).Where(`type = 2 AND created BETWEEN ? AND ? `, date+" 00:00:00", date+" 23:59:59").Scan(&diaries)
	for _, diary := range diaries {
		downloadPic(diary.PhotoUrl)
	}
}

func downloadPic(url string) {
	if url == "" {
		return
	}
	URL, _ := surl.Parse(url)
	suffixpath := URL.Path
	if strings.HasSuffix(URL.Path, "!large") {
		suffixpath = strings.TrimSuffix(URL.Path, "!large")
	} else if strings.Contains(URL.Path, "photos") {
		suffixpath = strings.Split(URL.Path, "photos")[1]
	}

	prepath := timepill.Conf.TimePill.PhotoPath + "/zhiding/" + path.Base(suffixpath)

	err := client.DownloadImage(prepath, url)
	if err != nil {
		log.Error(err)
	}
}
