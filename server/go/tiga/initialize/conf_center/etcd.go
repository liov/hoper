package conf_center

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Etcd struct {
	clientv3.Config
	Watch bool
	Key   string
}

// etcd
func (e *Etcd) HandleConfig(handle func([]byte)) error {
	client, err := clientv3.New(e.Config)
	if err != nil {
		return err
	}
	resp, err := client.Get(context.Background(), e.Key)
	if err != nil {
		return err
	}

	handle(resp.Kvs[0].Value)
	return nil
}
