package etcd

import (
	"context"
	initialize2 "github.com/liov/hoper/server/go/lib/initialize"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdConfig clientv3.Config

func (conf *EtcdConfig) Build() *clientv3.Client {
	client, _ := clientv3.New((clientv3.Config)(*conf))
	resp, _ := client.Get(context.Background(), initialize2.InitKey)
	initialize2.InitConfig.UnmarshalAndSet(resp.Kvs[0].Value)
	return client
}

func (conf *EtcdConfig) Generate() interface{} {
	return conf.Build()
}

type Eecd struct {
	*clientv3.Client
	Conf EtcdConfig
}

func (e *Eecd) Config() initialize2.Generate {
	return &e.Conf
}

func (e *Eecd) SetEntity(entity interface{}) {
	if client, ok := entity.(*clientv3.Client); ok {
		e.Client = client
	}
}

func (e *Eecd) Close() error {
	return e.Client.Close()
}
