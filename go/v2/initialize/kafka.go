package initialize

import (
	"github.com/Shopify/sarama"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/reflect3"
)

func (i *Init) P2Kafka() (sarama.SyncProducer,sarama.Consumer){
	conf := KafkaConfig{}
	if exist := reflect3.GetFieldValue(i.conf, &conf); !exist {
		return nil,nil
	}

	if conf.Model == 0 || conf.Model ==2{
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
		return producer,nil
	}
	if conf.Model == 1 || conf.Model ==2 {
		consumer, err := sarama.NewConsumer(conf.ConsAddr, nil)
		if err != nil {
			log.Info(err)
		}
		return nil,consumer
	}
	return nil,nil
}
