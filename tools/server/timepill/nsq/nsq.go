package nsq

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/nsq"
)

type Dao struct {
	ctx  context.Context
	NsqP *nsq.Producer
	NsqC *nsq.Consumer
}

func (dao *Dao) Send(topic string, msg []byte) error {
	return dao.NsqP.Publish(topic, msg)
}
