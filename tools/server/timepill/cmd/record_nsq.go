package main

import (
	"github.com/actliboy/hoper/server/go/lib/initialize"
	insq "github.com/actliboy/hoper/server/go/lib/initialize/nsq"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/nsqio/go-nsq"
	"tools/timepill"
)

type dao struct {
	NsqC  insq.Consumer `init:"config:nsq-consumer"`
	NsqC1 insq.Consumer `init:"config:nsq-consumer2"`
}

func (d *dao) Init() {
	d.NsqC.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Debug("消费者1收到消息：", string(message.Body))
		log.Debug(d.NsqC)
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
	if err := d.NsqC.ConnectToNSQDs(d.NsqC.Conf.NSQdAddrs); err != nil {
		log.Fatal(err)
	}
	d.NsqC1.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Debug("消费者2收到消息：", string(message.Body))
		log.Debug(d.NsqC1)
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
	if err := d.NsqC1.ConnectToNSQDs(d.NsqC1.Conf.NSQdAddrs); err != nil {
		log.Fatal(err)
	}
}

func main() {
	dao := &dao{}
	defer initialize.Start(&timepill.Conf, dao)()

	select {}
}
