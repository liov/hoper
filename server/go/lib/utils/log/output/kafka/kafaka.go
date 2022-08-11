package kafka

import (
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/url"
)

// kafka://${token}?sercret=${sercret}
func RegisterSink() {
	_ = zap.RegisterSink("kafka", func(url *url.URL) (sink zap.Sink, e error) {
		kl := new(LogKafka)
		kl.Topic = url.Query().Get("topic")
		// 设置日志输入到Kafka的配置
		config := sarama.NewConfig()
		//等待服务器所有副本都保存成功后的响应
		config.Producer.RequiredAcks = sarama.WaitForAll
		//随机的分区类型
		config.Producer.Partitioner = sarama.NewRandomPartitioner
		//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
		config.Producer.Return.Successes = true
		config.Producer.Return.Errors = true
		kl.Producer, e = sarama.NewSyncProducer([]string{url.Host}, config)
		return kl, nil
	})
}

type LogKafka struct {
	Producer sarama.SyncProducer
	Topic    string
}

func (lk *LogKafka) Write(b []byte) (n int, err error) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = lk.Topic
	msg.Value = sarama.ByteEncoder(b)
	_, _, err = lk.Producer.SendMessage(msg)
	if err != nil {
		return
	}
	return
}

func (lk *LogKafka) Sync() error {
	return nil
}

func (lk *LogKafka) Close() error {
	return lk.Producer.Close()
}

func NewCore() zapcore.Core { return core{} }

type core struct{}

func (core) Enabled(zapcore.Level) bool                                            { return false }
func (n core) With([]zapcore.Field) zapcore.Core                                   { return n }
func (core) Check(_ zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry { return ce }
func (core) Write(zapcore.Entry, []zapcore.Field) error                            { return nil }
func (core) Sync() error                                                           { return nil }
