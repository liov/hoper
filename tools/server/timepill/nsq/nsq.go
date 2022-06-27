package nsq

import (
	"context"
	insq "github.com/actliboy/hoper/server/go/lib/tiga/initialize/nsq"
)

type Dao struct {
	ctx  context.Context
	NsqP *insq.Producer
	NsqC *insq.Consumer
}

func (dao *Dao) Send(topic string, msg []byte) error {
	return dao.NsqP.Publish(topic, msg)
}
