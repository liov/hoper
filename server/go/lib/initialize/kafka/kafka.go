package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/actliboy/hoper/server/go/lib/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
)

type KafkaConfig struct {
	Addrs []string
	*sarama.Config
}

type KafkaProducerConfig KafkaConfig

func (conf *KafkaProducerConfig) generate() sarama.SyncProducer {

	// 等待服务器所有副本都保存成功后的响应
	conf.Producer.RequiredAcks = sarama.WaitForAll
	// 随机的分区类型：返回一个分区器，该分区器每次选择一个随机分区
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	// 是否等待成功和失败后的响应
	conf.Producer.Return.Successes = true

	// 使用给定代理地址和配置创建一个同步生产者
	producer, err := sarama.NewSyncProducer(conf.Addrs, conf.Config)
	if err != nil {
		log.Info(err)
	}
	//closes = append(closes,producer.CloseDao)
	return producer

}

func (conf *KafkaProducerConfig) Generate() interface{} {
	return conf.generate()
}

type KafkaConsumerConfig KafkaConfig

func (conf *KafkaConsumerConfig) generate() sarama.Consumer {

	consumer, err := sarama.NewConsumer(conf.Addrs, conf.Config)
	if err != nil {
		log.Info(err)
	}
	//closes = append(closes,consumer.CloseDao)
	return consumer

}

func (conf *KafkaConsumerConfig) Generate() interface{} {
	return conf.generate()
}

type KafkaProducer struct {
	sarama.SyncProducer
	Conf KafkaProducerConfig
}

func (k *KafkaProducer) Config() initialize.Generate {
	k.Conf.Config = sarama.NewConfig()
	k.Conf.Config.Version = sarama.V3_1_0_0
	return &k.Conf
}

func (k *KafkaProducer) SetEntity(entity interface{}) {
	if client, ok := entity.(sarama.SyncProducer); ok {
		k.SyncProducer = client
	}
}

type KafkaConsumer struct {
	sarama.Consumer
	Conf KafkaConsumerConfig
}

func (k *KafkaConsumer) Config() initialize.Generate {
	k.Conf.Config = sarama.NewConfig()
	k.Conf.Config.Version = sarama.V3_1_0_0
	return &k.Conf
}

func (k *KafkaConsumer) SetEntity(entity interface{}) {
	if client, ok := entity.(sarama.Consumer); ok {
		k.Consumer = client
	}
}
