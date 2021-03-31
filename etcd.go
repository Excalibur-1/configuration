package configuration

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

type etcdClient struct {
	c *clientv3.Client
}

func NewEtcdClient(endpoints []string, username, password string) (cl *etcdClient, err error) {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second * 5,
		Username:    username,
		Password:    password,
	})
	if err != nil {
		return
	}
	cl = &etcdClient{c: c}
	return
}

func (c *etcdClient) GetValues(keys []string) (vars map[string]string, err error) {
	vars = make(map[string]string)
	ctx := context.Background()
	for _, v := range keys {
		res, err := c.c.Get(ctx, v)
		if err != nil {
			return vars, err
		}
		kvs := res.Kvs
		for _, kv := range kvs {
			vars[v] = string(kv.Value)
		}
	}
	return
}

func (c *etcdClient) WatchPrefix(keys []string, waitIndex uint64) (index uint64, err error) {
	index = 1
	// return something > 0 to trigger a key retrieval from the store
	if waitIndex == 0 {
		return
	}

	respChan := make(chan WatchResponse)
	cancelRoutine := make(chan bool)
	defer close(cancelRoutine)

	// watch all keys in prefix for changes
	for _, v := range keys {
		log.Printf("Watching: %v\n", v)
		go c.watch(v, respChan, cancelRoutine)
	}

	for {
		select {
		case r := <-respChan:
			index = r.WaitIndex
			err = r.Err
			return
		}
	}
}

func (c *etcdClient) watch(key string, respChan chan WatchResponse, cancelRoutine chan bool) {
	ctx := context.Background()
	watchCh := c.c.Watch(ctx, key)
	for {
		select {
		case <-cancelRoutine:
			log.Printf("Stop watching:%v\n", key)
			return
		default:
			for resp := range watchCh {
				events := resp.Events
				for range events {
					respChan <- WatchResponse{WaitIndex: 1}
				}
			}
		}
	}
}
