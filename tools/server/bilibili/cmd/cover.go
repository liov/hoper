package main

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/initialize"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"tools/bilibili/config"
	"tools/bilibili/dao"
	"tools/bilibili/download"
	"tools/bilibili/rpc"
	"tools/bilibili/tool"
)

func main() {
	defer initialize.Start(config.Conf, &dao.Dao)()

	config.Conf.Bilibili.DownloadPicPath = "F:\\B站\\pic"
	apiservice := rpc.API{}
	page := 1
	for {
		log.Printf("第%d页\n", page)
		res, err := apiservice.GetFavLResourceList(63181530, page)
		if err != nil {
			log.Println(err)
		}

		for _, fav := range res.Medias {
			aid := tool.Bv2av(fav.Bvid)
			filepath := filepath.Join(config.Conf.Bilibili.DownloadPicPath, strconv.Itoa(fav.Upper.Mid), strconv.Itoa(fav.Upper.Mid)+"_"+strconv.Itoa(aid)+"_"+path.Base(fav.Cover))
			_, err := os.Stat(filepath)
			if err == nil {
				return
			}
			err = download.CoverDownload(context.Background(), fav.Cover, fav.Upper.Mid, aid)
			if err != nil {
				log.Println("下载图片失败：", err)
			}
		}
		if len(res.Medias) == 0 {
			return
		}
		page++
	}

}
