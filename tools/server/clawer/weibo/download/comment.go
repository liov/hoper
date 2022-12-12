package download

import (
	"context"
	"github.com/liov/hoper/server/go/lib/utils/generics/net/http/client/crawler"
	timei "github.com/liov/hoper/server/go/lib/utils/time"
	"strconv"
	"tools/clawer/weibo/dao"
	"tools/clawer/weibo/rpc"
)

func GetCommentReq(wid string, page int) *crawler.Request {
	return crawler.NewRequest("GetCommentReq"+wid+" "+strconv.Itoa(page), KindGet, func(ctx context.Context) ([]*crawler.Request, error) {
		commentList, err := rpc.GetComments(wid, page)
		if err != nil {
			return nil, err
		}
		var comments []*dao.Comment
		if commentList != nil {
			for _, comment := range commentList.Data {
				createdAt := timei.WeiboParse(comment.CreatedAt)
				comments = append(comments, &dao.Comment{
					Id:        comment.Id,
					UserId:    comment.User.Id,
					CreatedAt: createdAt,
					Source:    comment.Source,
					Text:      comment.Text,
					ReplyId:   comment.ReplyId,
					ReplyText: comment.ReplyText,
				})
			}
			dao.Dao.Hoper.Create(comments)
		}
		if commentList != nil && len(commentList.Data) > 0 {
			return []*crawler.Request{GetCommentReq(wid, page+1)}, nil
		}
		return nil, nil
	})
}
