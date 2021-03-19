package initialize

import (
	"github.com/Shopify/sarama"
	"github.com/liov/hoper/go/v2/utils/log"
)

type KafkaConfig struct {
	Model    int8 //0生产者，1消费者，2所有
	Topic    string
	ProdAddr []string
	ConsAddr []string
}

func (conf *KafkaConfig) Generate() (sarama.SyncProducer, sarama.Consumer) {
	if conf.Model == 0 || conf.Model == 2 {
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
		return producer, nil
	}
	if conf.Model == 1 || conf.Model == 2 {
		consumer, err := sarama.NewConsumer(conf.ConsAddr, nil)
		if err != nil {
			log.Info(err)
		}
		//closes = append(closes,consumer.CloseDao)
		return nil, consumer
	}
	return nil, nil
}

func (init *Inject) P2Kafka() (sarama.SyncProducer, sarama.Consumer) {
	conf := &KafkaConfig{}
	if exist := init.SetConf(conf); !exist {
		return nil,nil
	}
	return conf.Generate()
}
