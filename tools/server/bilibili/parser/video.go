package parser

import (
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"github.com/tidwall/gjson"
	"tools/bilibili/fetcher"
	"tools/bilibili/model"
)

// 大视频所以分成了不同部分提交，但是最终显示的只有一个视频文件
func GenVideoDownloadParseFun(videoCid *model.VideoCid) crawler.ParseFun {
	return func(contents []byte) ([]*crawler.Request, error) {
		var requests []*crawler.Request

		durlSlice := gjson.GetBytes(contents, "durl").Array()
		videoCid.AllOrder = int64(len(durlSlice))

		for _, i := range durlSlice {
			video := &model.Video{Order: i.Get("order").Int(), ParCid: videoCid}
			videoUrl := i.Get("url").String()
			req := crawler.NewRequest2(videoUrl, crawler.FetchFun(fetcher.GenVideoFetcher(video)), recordCidParseFun(video))
			requests = append(requests, req)
		}
		return requests, nil
	}
}

func recordCidParseFun(Video *model.Video) crawler.ParseFun {
	return func(contents []byte) ([]*crawler.Request, error) {
		return nil, nil
	}
}
