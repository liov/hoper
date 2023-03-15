package nsq

import (
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

func (p *Producer) Config() any {
	p.Conf.Config = nsq.NewConfig()
	return &p.Conf
}

func (p *Producer) SetEntity(entity interface{}) {
	p.Producer = p.Conf.Build()
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

type Consumer struct {
	*nsq.Consumer `init:"entity"`
	Conf          ConsumerConfig `init:"config"`
}

func (c *Consumer) Config() any {
	c.Conf.Config = nsq.NewConfig()
	return &c.Conf
}

func (c *Consumer) SetEntity() {
	c.Consumer = c.Conf.Build()
}

func (c *Consumer) Close() error {
	c.Consumer.Stop()
	return nil
}
