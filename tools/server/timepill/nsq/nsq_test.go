package nsq

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	insq "github.com/actliboy/hoper/server/go/lib/tiga/initialize/nsq"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/nsqio/go-nsq"
	"testing"
	"tools/timepill"
)

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
	d := &dao{}
	defer initialize.Start(&timepill.Conf, d)()
	t.Log(d.NsqP.Publish("diary", []byte("test")))
}

func TestNsqc(t *testing.T) {
	d := &dao{}
	defer initialize.Start(&timepill.Conf, d)()
	/*	if err := d.NsqC.ConnectToNSQLookupds(d.NsqC.Conf.NSQLookupdAddrs); err != nil {
		log.Fatal(err)
	}*/
	if err := d.NsqC.ConnectToNSQDs(d.NsqC.Conf.NSQdAddrs); err != nil {
		log.Fatal(err)
	}
	t.Log("nsq test")
	select {}
}
