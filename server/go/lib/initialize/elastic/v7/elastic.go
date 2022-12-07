package v7

import (
	"github.com/liov/hoper/server/go/lib/initialize"
	"github.com/liov/hoper/server/go/lib/utils/log"
	"github.com/olivere/elastic/v7"
	"github.com/olivere/elastic/v7/config"
)

type ElasticConfig config.Config

func (conf *ElasticConfig) Build() *elastic.Client {
	client, err := elastic.NewClientFromConfig((*config.Config)(conf))
	if err != nil {
		log.Error(err)
	}
	//closes = append(closes,client.Stop)
	return client
}

func (conf *ElasticConfig) Generate() interface{} {
	return conf.Build()
}

type Es struct {
	*elastic.Client
	Conf ElasticConfig
}

func (es *Es) Config() initialize.Generate {
	return &es.Conf
}

func (es *Es) SetEntity(entity interface{}) {
	if client, ok := entity.(*elastic.Client); ok {
		es.Client = client
	}
}

func (es *Es) Close() error {
	return nil
}
