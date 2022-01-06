package inject_dao

import (
	"github.com/Shopify/sarama"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
)

type KafkaProducerConfig struct {
	Topic    string
	ProdAddr []string
}

func (conf *KafkaProducerConfig) generate() sarama.SyncProducer {

	config := sarama.NewConfig()
	// 等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 随机的分区类型：返回一个分区器，该分区器每次选择一个随机分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 是否等待成功和失败后的响应
	config.Producer.Return.Successes = true

	// 使用给定代理地址和配置创建一个同步生产者
	producer, err := sarama.NewSyncProducer(conf.ProdAddr, config)
	if err != nil {
		log.Info(err)
	}
	//closes = append(closes,producer.CloseDao)
	return producer

}

func (conf *KafkaProducerConfig) Generate() interface{} {
	return conf.generate()
}

type KafkaConsumerConfig struct {
	Topic    string
	ConsAddr []string
}

func (conf *KafkaConsumerConfig) generate() sarama.Consumer {

	consumer, err := sarama.NewConsumer(conf.ConsAddr, nil)
	if err != nil {
		log.Info(err)
	}
	//closes = append(closes,consumer.CloseDao)
	return consumer

}

func (conf *KafkaConsumerConfig) Generate() interface{} {
	return conf.generate()
}
