package initialize

import (
	"log"

	"github.com/nsqio/go-nsq"
)

type NsqConfig struct {
	Addr    string
	Model   int8 //0生产者，1消费者，2所有
	Topic   string
	Channel string
}

func (conf *NsqConfig) Generate() (*nsq.Producer, *nsq.Consumer) {
	cfg := nsq.NewConfig()
	if conf.Model == 0 || conf.Model == 2 {
		producer, err := nsq.NewProducer(conf.Addr, cfg)
		if err != nil {
			panic(err)
		}
		//closes = append(closes,producer.Stop)
		return producer, nil
	}
	if conf.Model == 1 || conf.Model == 2 {
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
		return nil, customer
	}
	return nil, nil
}

func (init *Inject) P2NSQ() (*nsq.Producer, *nsq.Consumer) {
	conf := &NsqConfig{}
	if exist := init.SetConf(conf); !exist {
		return nil,nil
	}
	return conf.Generate()
}
