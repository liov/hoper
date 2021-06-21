package inject_dao

import (
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdConfig clientv3.Config

func (conf *EtcdConfig) generate() *clientv3.Client {
	client, _ := clientv3.New((clientv3.Config)(*conf))

	return client
}

func (conf *EtcdConfig) Generate() interface{} {
	return conf.generate()
}
