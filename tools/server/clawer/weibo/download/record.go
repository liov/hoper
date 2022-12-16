package download

import (
	"context"
	"encoding/json"
	"github.com/liov/hoper/server/go/lib/utils/log"
	stringsi "github.com/liov/hoper/server/go/lib/utils/strings"
	"github.com/liov/hoper/server/go/lib_v2/utils/net/http/client/crawler"
	"strconv"
	"strings"
	"time"
	"tools/clawer/weibo/dao"
	"tools/clawer/weibo/rpc"
)

func RecordUsersWeiboReq(userIds []int, record bool) []*crawler.Request {
	var reqs []*crawler.Request
	for _, userId := range userIds {
		reqs = append(reqs, RecordUserWeiboReq(userId, 1, record))
	}
	return reqs
}

func RecordUserWeiboReq(uid, page int, record bool) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{BaseTaskMeta: crawler.BaseTaskMeta{Key: strconv.Itoa(uid) + " " + strconv.Itoa(page) + "RecordUserWeiboReq"}, Kind: KindGet},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			log.Infof("RecordUserWeiboReq %d 第%d页", uid, page)
			piccards, err := rpc.GetChannels(rpc.ALL, uid, page)
			var requests []*crawler.Request
			if err != nil {
				if strings.HasPrefix(err.Error(), "status:403") {

				}
				if _, ok := err.(*json.UnmarshalTypeError); ok {
					ppiccards, err := rpc.GetChannelsV2(rpc.ALL, uid, page)
					if err != nil {
						return nil, err
					}
					for _, card := range ppiccards.Cards {

						if card.Mblog != nil {
							if record {
								exists, err := dao.WeiboExists(card.Mblog.Id)
								if err != nil {
									continue
								}
								if exists {
									return requests, nil
								}
							}
							requests = append(requests, LongWeiboReq(card.Mblog.Mid, record))
						}
					}
					/*	if len(ppiccards.Cards) > 2 {
						requests = append(requests, RecordUserWeiboReq(uid, page+1, record))
					}*/
					return requests, err
				}
				return nil, err
			}

			if piccards.Cards != nil {
				for _, card := range piccards.Cards {
					if card.Mblog != nil {
						if record {
							exists, err := dao.WeiboExists(card.Mblog.Id)
							if err != nil {
								continue
							}
							if exists {
								return requests, nil
							}
						}

						if card.Mblog.PicNum > 9 || card.Mblog.IsLongText {
							requests = append(requests, LongWeiboReq(card.Mblog.Id, record))
						} else {
							requests = append(requests, GetWeiboReq(card.Mblog, record)...)
						}
					}

				}
			}
			/*	if len(piccards.Cards) > 2 {
				requests = append(requests, RecordUserWeiboReq(uid, page+1, record))
			}*/
			return requests, nil
		},
	}
}

func GetWeiboReq(mblog *rpc.Mblog, record bool) []*crawler.Request {
	var requests []*crawler.Request
	id, _ := strconv.Atoi(mblog.Id)
	createdAt, _ := time.Parse(time.RubyDate, mblog.CreatedAt)
	var video string
	if mblog.PageInfo != nil && mblog.PageInfo.Type == "video" {
		requests = append(requests, DownloadVideoReq(mblog))
		var url string
		if mblog.PageInfo.Urls.Mp4720PMp4 != "" {
			url = mblog.PageInfo.Urls.Mp4720PMp4
		} else if mblog.PageInfo.Urls.Mp4HdMp4 != "" {
			url = mblog.PageInfo.Urls.Mp4HdMp4
		} else {
			url = mblog.PageInfo.Urls.Mp4LdMp4
		}
		if url != "" {
			video = stringsi.CountdownCutoff(stringsi.CutoffContain(url, "mp4"), "/")
		}
	}

	exists, _ := dao.UserExists(mblog.User.Id)
	if !exists {
		dao.Dao.Hoper.Create(&dao.User{
			Id:          mblog.User.Id,
			ScreenName:  mblog.User.ScreenName,
			Gender:      mblog.User.Gender,
			Description: mblog.User.Description,
			AvatarHd:    mblog.User.AvatarHd,
		})
	}

	if record {
		var retweetedId int
		if mblog.RetweetedStatus != nil {
			retweetedId, _ = strconv.Atoi(mblog.RetweetedStatus.Id)
		}
		dao.Dao.Hoper.Create(&dao.Weibo{
			Id:          id,
			BId:         mblog.Bid,
			UserId:      mblog.User.Id,
			Text:        mblog.Text,
			CreatedAt:   createdAt,
			Source:      mblog.Source,
			Pics:        strings.Join(mblog.PicIds, ","),
			Video:       video,
			RetweetedId: retweetedId,
		})
		if mblog.CommentsCount > 0 {
			requests = append(requests, GetCommentReq(mblog.Id, 1))
		}
	}

	requests = append(requests, DownloadPhotoReqsV2(mblog)...)
	return requests
}
