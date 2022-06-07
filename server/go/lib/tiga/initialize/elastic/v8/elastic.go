package v8

import (
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/auth"
	"github.com/elastic/go-elasticsearch/v8"
	"net/http"
)

type ElasticConfig elasticsearch.Config

func (conf *ElasticConfig) generate() *elasticsearch.Client {
	conf.Header = http.Header{}
	if conf.Username != "" && conf.Password != "" {
		auth.SetBasicAuth(conf.Header, conf.Username, conf.Password)
	}
	client, err := elasticsearch.NewClient((elasticsearch.Config)(*conf))
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
	*elasticsearch.Client
	Conf ElasticConfig
}

func (es *Es) Config() interface{} {
	return &es.Conf
}

func (es *Es) SetEntity(entity interface{}) {
	if client, ok := entity.(*elasticsearch.Client); ok {
		es.Client = client
	}
}
