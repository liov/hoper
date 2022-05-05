package elastic

import (
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/olivere/elastic/v7"
	"github.com/olivere/elastic/v7/config"
)

type ElasticConfig config.Config

func (conf *ElasticConfig) generate() *elastic.Client {
	client, err := elastic.NewClientFromConfig((*config.Config)(conf))
	if err != nil {
		log.Error(err)
	}
	//closes = append(closes,client.Stop)
	return client
}

func (conf *ElasticConfig) Generate() interface{} {
	return conf.generate()
}

type Es struct {
	*elastic.Client
	Conf ElasticConfig
}

func (es *Es) Config() interface{} {
	return &es.Conf
}

func (es *Es) SetEntity(entity interface{}) {
	if client, ok := entity.(*elastic.Client); ok {
		es.Client = client
	}
}
