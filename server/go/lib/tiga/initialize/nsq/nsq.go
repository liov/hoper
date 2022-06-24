package nsq

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"log"

	"github.com/nsqio/go-nsq"
)

type ProducerConfig struct {
	Addr string
	*nsq.Config
}

func (conf *ProducerConfig) generate() *nsq.Producer {

	producer, err := nsq.NewProducer(conf.Addr, conf.Config)
	if err != nil {
		panic(err)
	}
	//closes = append(closes,producer.Stop)
	return producer
}

func (conf *ProducerConfig) Generate() interface{} {
	return conf.generate()
}

type ConsumerConfig struct {
	Addr    string
	Topic   string
	Channel string
	*nsq.Config
}

func (conf *ConsumerConfig) generate() *nsq.Consumer {

	customer, err := nsq.NewConsumer(conf.Topic, conf.Channel, conf.Config)
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

func (conf *ConsumerConfig) Generate() interface{} {
	return conf.generate()
}

type Producer struct {
	*nsq.Producer
	Conf ProducerConfig
}

func (p *Producer) Config() initialize.Generate {
	p.Conf.Config = nsq.NewConfig()
	return &p.Conf
}

func (p *Producer) SetEntity(entity interface{}) {
	if client, ok := entity.(*nsq.Producer); ok {
		p.Producer = client
	}
}

type Consumer struct {
	*nsq.Consumer
	Conf ProducerConfig
}

func (c *Consumer) Config() interface{} {
	c.Conf.Config = nsq.NewConfig()
	return &c.Conf
}

func (c *Consumer) SetEntity(entity interface{}) {
	if client, ok := entity.(*nsq.Consumer); ok {
		c.Consumer = client
	}
}
