package initialize

import (
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/reflect3"
	"github.com/olivere/elastic"
)

func (i *Init) P2Elastic() {
	esConf:=ElasticConfig{}
	if exist := reflect3.GetFieldValue(i.conf,&esConf);!exist{
		return
	}
	client,err:=elastic.NewClient()
	if err!=nil{
		log.Error(err)
	}
	reflect3.SetFieldValue(i.dao,client)
}

