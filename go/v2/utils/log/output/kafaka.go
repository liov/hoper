package output

import "github.com/Shopify/sarama"

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
