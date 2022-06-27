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

type Producer struct {
	*nsq.Producer `init:"entity"`
	Conf          ProducerConfig `init:"config"`
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
	NSQLookupdAddrs []string
	NSQdAddrs       []string
	Topic           string
	Channel         string
	*nsq.Config
}

func (conf *ConsumerConfig) generate() *nsq.Consumer {

	customer, err := nsq.NewConsumer(conf.Topic, conf.Channel, conf.Config)
	if err != nil {
		log.Fatal(err)
	}

	/*	if err := customer.ConnectToNSQLookupds(conf.NSQLookupdAddrs); err != nil {
		log.Fatal(err)
	}*/

	/*	if err = customer.ConnectToNSQDs(conf.NSQdAddrs); err != nil {
			log.Fatal(err)
		}
	*/
	return customer

}

func (conf *ConsumerConfig) Generate() interface{} {
	return conf.generate()
}

type Consumer struct {
	*nsq.Consumer `init:"entity"`
	Conf          ConsumerConfig `init:"config"`
}

func (c *Consumer) Config() initialize.Generate {
	c.Conf.Config = nsq.NewConfig()
	return &c.Conf
}

func (c *Consumer) SetEntity(entity interface{}) {
	if client, ok := entity.(*nsq.Consumer); ok {
		c.Consumer = client
	}
}
