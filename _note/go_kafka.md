```go
	// 等待服务器所有副本都保存成功后的响应
conf.Producer.RequiredAcks = sarama.WaitForAll
// 随机的分区类型：返回一个分区器，该分区器每次选择一个随机分区
conf.Producer.Partitioner = sarama.NewRandomPartitioner
// 是否等待成功和失败后的响应
conf.Producer.Return.Successes = true
```