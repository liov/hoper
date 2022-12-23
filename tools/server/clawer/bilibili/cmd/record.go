package main

import (
	"context"
	"github.com/liov/hoper/server/go/lib/initialize"
	"log"
	"tools/clawer/bilibili/config"
	"tools/clawer/bilibili/dao"
	"tools/clawer/bilibili/download"
)

func main() {
	defer initialize.Start(config.Conf, &dao.Dao)()
	//delete("F:\\B站\\video\\10139490\\10139490_207568591_395475557_～Alone～_alone_1_120.flv")
	//deduplication()
	ctx := context.Background()
	view2, err := download.RecordViewInfo(ctx, 800947366)
	if view2 == nil {
		log.Println(err)
	}
	for _, page := range view2.Pages {
		if len(view2.Pages) == 1 {
			page.Part = download.PartEqTitle
		}
		video := download.NewVideo(view2.Owner.Mid, view2.Title, view2.Aid, page.Cid, page.Page, page.Part, view2.PubDate)
		_, err := video.RecordVideo(ctx)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
