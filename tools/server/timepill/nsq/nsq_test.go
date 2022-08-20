package nsq

import (
	"github.com/actliboy/hoper/server/go/lib/initialize"
	insq "github.com/actliboy/hoper/server/go/lib/initialize/nsq"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/nsqio/go-nsq"
	"testing"
)

type config struct{}

func (c *config) Init() {
}

type dao struct {
	NsqP insq.Producer `init:"config:nsq-producer"`
	NsqC insq.Consumer `init:"config:nsq-consumer"`
}

func (d *dao) Init() {
	d.NsqC.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Info(string(message.Body))
		return nil
	}))
}

func (d *dao) Close() {

}

func TestNsqp(t *testing.T) {
	c := &config{}
	d := &dao{}
	defer initialize.Start(c, d)()
	t.Log(d.NsqP.Publish("timepill_pic", []byte("test")))
}

func TestNsqc(t *testing.T) {
	c := &config{}
	d := &dao{}
	defer initialize.Start(c, d)()
	/*	if err := d.NsqC.ConnectToNSQLookupds(d.NsqC.Conf.NSQLookupdAddrs); err != nil {
		log.Fatal(err)
	}*/
	if err := d.NsqC.ConnectToNSQDs(d.NsqC.Conf.NSQdAddrs); err != nil {
		log.Fatal(err)
	}
	t.Log("nsq test")
	select {}
}
