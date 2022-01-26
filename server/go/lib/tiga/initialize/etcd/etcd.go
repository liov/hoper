package etcd

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdConfig clientv3.Config

func (conf *EtcdConfig) generate() *clientv3.Client {
	client, _ := clientv3.New((clientv3.Config)(*conf))
	resp, _ := client.Get(context.Background(), initialize.InitKey)
	initialize.InitConfig.UnmarshalAndSet(resp.Kvs[0].Value)
	return client
}

func (conf *EtcdConfig) Generate() interface{} {
	return conf.generate()
}
