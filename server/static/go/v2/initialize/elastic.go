package initialize

import (
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/reflect3"
	"github.com/olivere/elastic"
)

type ElasticConfig struct {
	Host string
	Port int32
}

func (conf *ElasticConfig) Generate() *elastic.Client {
	client, err := elastic.NewClient()
	if err != nil {
		log.Error(err)
	}
	//closes = append(closes,client.Stop)
	return client
}

func (init *Init) P2Elastic() *elastic.Client {
	conf := &ElasticConfig{}
	if exist := reflect3.GetFieldValue(init.conf, conf); !exist {
		return nil
	}
	return conf.Generate()
}
