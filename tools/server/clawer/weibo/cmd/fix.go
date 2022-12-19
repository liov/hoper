package main

import (
	"github.com/liov/hoper/server/go/lib/initialize"
	"github.com/liov/hoper/server/go/lib/utils/fs"
	stringsi "github.com/liov/hoper/server/go/lib/utils/strings"
	"github.com/liov/hoper/server/go/lib_v2/utils/net/http/client/crawler"
	"log"
	"os"
	"strings"
	"time"
	claweri "tools/clawer"
	"tools/clawer/weibo/config"
	"tools/clawer/weibo/dao"
	"tools/clawer/weibo/download"
)

func main() {
	defer initialize.Start(config.Conf, &dao.Dao)()
	dao.Dao.Hoper.Migrator().CreateTable(&claweri.Dir{})
	engine := crawler.NewEngine(config.Conf.Weibo.WorkCount).Timer(download.KindGet, time.Second)

	engine.Run(download.GetUserFollowWeiboReq(""))
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
