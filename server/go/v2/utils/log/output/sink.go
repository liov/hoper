package output

import (
	"net/url"
	"os"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type nopCloserSink struct{ zapcore.WriteSyncer }

func (nopCloserSink) Close() error { return nil }

func RegisterSink() {
	_ = zap.RegisterSink("http", func(url *url.URL) (sink zap.Sink, e error) {
		call := new(LoggerCall)
		call.HookURL = []string{url.String()}
		return nopCloserSink{zapcore.AddSync(call)}, nil
	})
	_ = zap.RegisterSink("https", func(url *url.URL) (sink zap.Sink, e error) {
		call := new(LoggerCall)
		call.HookURL = []string{url.String()}
		return nopCloserSink{zapcore.AddSync(call)}, nil
	})
	_ = zap.RegisterSink("socket", func(url *url.URL) (sink zap.Sink, e error) {
		return nopCloserSink{zapcore.AddSync(os.Stdout)}, nil
	})

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
		return nopCloserSink{zapcore.AddSync(kl)}, nil
	})
}
