package v7

import (
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

type Es struct {
	*elastic.Client
	Conf ElasticConfig
}

func (es *Es) Config() any {
	return &es.Conf
}

func (es *Es) SetEntity() {
	es.Client = es.Conf.Build()
}

func (es *Es) Close() error {
	return nil
}
