package initialize

import (
	"log"

	"github.com/liov/hoper/go/v2/utils/reflect3"
	"github.com/nsqio/go-nsq"
)

func (i *Init) P2NSQ() {
	nsqConf := NsqConfig{}
	if exist := reflect3.GetFieldValue(i.conf, &nsqConf); !exist {
		return
	}
	cfg := nsq.NewConfig()
	if nsqConf.model == 0 || nsqConf.model == 2 {
		producer, err := nsq.NewProducer(nsqConf.Host, cfg)
		if err != nil {
			panic(err)
		}
		reflect3.SetFieldValue(i.dao,producer)
	}
	if nsqConf.model == 1 || nsqConf.model == 2 {
		customer, err := nsq.NewConsumer(nsqConf.topic, nsqConf.channel, cfg)
		if err != nil {
			panic(err)
		}
		//c.SetLogger(nil, 0)       //屏蔽系统日志
		//customer.AddHandler(handle) // 添加消费者接口

		//建立NSQLookupd连接
		if err := customer.ConnectToNSQLookupd(nsqConf.Host); err != nil {
			log.Println("consumer 新建失败")
		}
		//建立多个nsqd连接
		// if err := customer.ConnectToNSQDs([]string{"127.0.0.1:4150", "127.0.0.1:4152"}); err != nil {
		//  panic(err)
		// }

		// 建立一个nsqd连接
		/*	if err := customer.ConnectToNSQD(address); err != nil {
			 panic(err)
			}
			<-c.StopChan*/
		reflect3.SetFieldValue(i.dao,customer)
	}
}
