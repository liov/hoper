package nsq

import (
	"context"
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"tools/timepill/model"
)

const (
	TopicDiary = "diary"
	TopicPic   = "timepill_pic"
	TopicCover = "timepill_cover"
)

type Dao struct {
	ctx  context.Context
	NsqP *nsq.Producer
	NsqC *nsq.Consumer
}

func (dao *Dao) Send(topic string, msg []byte) error {
	return dao.NsqP.Publish(topic, msg)
}

func PublishPic(nsqp *nsq.Producer, userId int, photoUrl, created string) error {
	data, err := json.Marshal(&Diary{
		UserId:   userId,
		PhotoUrl: photoUrl,
		Created:  created,
	})
	if err != nil {
		return err
	}
	return nsqp.Publish(TopicPic, data)
}

func PublishCover(nsqp *nsq.Producer, typ model.CoverType, url string) error {
	data, err := json.Marshal(&Cover{
		typ, url,
	})
	if err != nil {
		return err
	}
	return nsqp.Publish(TopicCover, data)
}
