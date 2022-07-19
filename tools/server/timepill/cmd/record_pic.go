package main

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	insq "github.com/actliboy/hoper/server/go/lib/tiga/initialize/nsq"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/nsqio/go-nsq"
	"tools/timepill"
)

type dao struct {
	initialize.DaoPlaceholder
	NsqC  insq.Consumer `init:"config:nsq-consumer"`
	NsqC1 insq.Consumer `init:"config:nsq-consumer2"`
}

func main() {
	dao := &dao{}
	defer initialize.Start(&timepill.Conf, dao)()
	dao.NsqC.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Debug("消费者1收到消息：", string(message.Body))
		log.Debug(dao.NsqC)
		/*diary := &tnsq.Diary{}
		err := json.Unmarshal(message.Body, diary)
		if err != nil {
			return err
		}
		err = timepill.DownloadPic(diary.UserId, diary.PhotoUrl, diary.Created)
		if err != nil {
			log.Error(err)
		}*/
		return nil
	}))
	if err := dao.NsqC.ConnectToNSQDs(dao.NsqC.Conf.NSQdAddrs); err != nil {
		log.Fatal(err)
	}
	dao.NsqC1.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Debug("消费者2收到消息：", string(message.Body))
		log.Debug(dao.NsqC1)
		/*cover := &tnsq.Cover{}
		err := json.Unmarshal(message.Body, cover)
		if err != nil {
			return err
		}
		err = timepill.DownloadCover(cover.Type.String(), cover.Url)
		if err != nil {
			log.Error(err)
		}*/
		return nil
	}))
	if err := dao.NsqC1.ConnectToNSQDs(dao.NsqC1.Conf.NSQdAddrs); err != nil {
		log.Fatal(err)
	}
	select {}
}
