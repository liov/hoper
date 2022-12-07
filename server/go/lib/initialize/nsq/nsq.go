package nsq

import (
	"github.com/liov/hoper/server/go/lib/initialize"
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

func (p *Producer) Close() error {
	p.Producer.Stop()
	return nil
}

func (conf *ProducerConfig) Build() *nsq.Producer {

	producer, err := nsq.NewProducer(conf.Addr, conf.Config)
	if err != nil {
		panic(err)
	}
	//closes = append(closes,producer.Stop)
	return producer
}

func (conf *ProducerConfig) Generate() interface{} {
	return conf.Build()
}

type ConsumerConfig struct {
	NSQLookupdAddrs []string
	NSQdAddrs       []string
	Topic           string
	Channel         string
	*nsq.Config
}

func (conf *ConsumerConfig) Build() *nsq.Consumer {

	consumer, err := nsq.NewConsumer(conf.Topic, conf.Channel, conf.Config)
	if err != nil {
		log.Fatal(err)
	}

	/*	if len(conf.NSQLookupdAddrs) > 0 {
			if err := consumer.ConnectToNSQLookupds(conf.NSQLookupdAddrs); err != nil {
				log.Fatal(err)
			}
		}
		if len(conf.NSQdAddrs) > 0 {
			if err = consumer.ConnectToNSQDs(conf.NSQdAddrs); err != nil {
				log.Fatal(err)
			}

		}*/
	return consumer

}

func (conf *ConsumerConfig) Generate() interface{} {
	return conf.Build()
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

func (c *Consumer) Close() error {
	c.Consumer.Stop()
	return nil
}
