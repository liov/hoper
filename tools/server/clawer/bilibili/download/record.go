package download

import (
	dbi "github.com/liov/hoper/server/go/lib/utils/dao/db/const"
	"github.com/liov/hoper/server/go/lib_v2/utils/net/http/client/crawler"
	"gorm.io/gorm"
	"time"
	"tools/clawer/bilibili/dao"
)

func DownloadRecordVideo(engine *crawler.Engine) {
	now := time.Now()
	for {
		gorm.Model{}
		var videos []*Video
		dao.Dao.Hoper.DB.Table(dao.TableNameVideo+" a").Select(`b.uid,b.aid,b.title,a.cid,a.page,a.part,a.record,b.pubdate pub_at,a.created_at`).Joins(`LEFT JOIN `+dao.TableNameView+" b ON a.aid = b.aid").Where(`b.record < 2  AND b.created_at < ? AND b.`+dbi.NotDeleted+`AND a.`+dbi.NotDeleted, now).Order(`b.created_at DESC`).Limit(100).Scan(&videos)
		if len(videos) == 0 {
			return
		}
		for _, video := range videos {
			if video.Title == "" {
				req := ViewRecordUpdateReqAfterRecordVideo(video.Aid)
				engine.BaseEngine.AddTask(engine.BaseTask(req))
			} else {
				req := video.GetVideoReqAfterDownloadVideo()
				engine.BaseEngine.AddTask(engine.BaseTask(req))
			}
		}
		now = videos[len(videos)-1].CreatedAt
	}
}
