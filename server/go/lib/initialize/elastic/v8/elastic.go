package v8

import (
	"github.com/actliboy/hoper/server/go/lib/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/auth"
	"github.com/elastic/go-elasticsearch/v8"
	"net/http"
)

type ElasticConfig elasticsearch.Config

func (conf *ElasticConfig) Build() *elasticsearch.Client {
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
	return conf.Build()
}

type Es struct {
	*elasticsearch.Client
	Conf ElasticConfig
}

func (es *Es) Config() initialize.Generate {
	return &es.Conf
}

func (es *Es) SetEntity(entity interface{}) {
	if client, ok := entity.(*elasticsearch.Client); ok {
		es.Client = client
	}
}

func (es *Es) Close() {
}

type Esv2 struct {
	*elasticsearch.Client `init:"entity"`
	Config                ElasticConfig `init:"config"`
}

func (es *Esv2) Close() {
}
