package etcd

import (
	"context"
	"github.com/liov/hoper/server/go/lib/initialize"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdConfig clientv3.Config

func (conf *EtcdConfig) Build() *clientv3.Client {
	client, _ := clientv3.New((clientv3.Config)(*conf))
	resp, _ := client.Get(context.Background(), initialize.InitKey)
	initialize.InitConfig.UnmarshalAndSet(resp.Kvs[0].Value)
	return client
}

type Eecd struct {
	*clientv3.Client
	Conf EtcdConfig
}

func (e *Eecd) Config() any {
	return &e.Conf
}

func (e *Eecd) SetEntity() {
	e.Client = e.Conf.Build()
}

func (e *Eecd) Close() error {
	return e.Client.Close()
}
