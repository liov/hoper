package inject_dao

import (
	"log"

	"github.com/nsqio/go-nsq"
)

type NsqProducerConfig struct {
	Addr    string
	Topic   string
	Channel string
}

func (conf *NsqProducerConfig) generate() *nsq.Producer {
	cfg := nsq.NewConfig()

	producer, err := nsq.NewProducer(conf.Addr, cfg)
	if err != nil {
		panic(err)
	}
	//closes = append(closes,producer.Stop)
	return producer
}

func (conf *NsqProducerConfig) Generate() interface{} {
	return conf.generate()
}

type NsqConsumerConfig struct {
	Addr    string
	Model   int8 //0生产者，1消费者，2所有
	Topic   string
	Channel string
}

func (conf *NsqConsumerConfig) generate() *nsq.Consumer {
	cfg := nsq.NewConfig()

	customer, err := nsq.NewConsumer(conf.Topic, conf.Channel, cfg)
	if err != nil {
		panic(err)
	}
	//c.SetLogger(nil, 0)       //屏蔽系统日志
	//customer.AddHandler(handle) // 添加消费者接口

	//建立NSQLookupd连接
	if err := customer.ConnectToNSQLookupd(conf.Addr); err != nil {
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
	//closes = append(closes,customer.Stop)
	return customer

}

func (conf *NsqConsumerConfig) Generate() interface{} {
	return conf.generate()
}
