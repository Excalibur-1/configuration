package configuration

import (
	"log"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

const DefaultSessionTimeout = time.Second * 10

// Client provides a wrapper around the zookeeper client
type zkClient struct {
	c *zk.Conn
}

func NewZkClient(machines []string, user, password, openUser, openPassword string) (cl *zkClient, err error) {
	c, _, err := zk.Connect(machines, DefaultSessionTimeout)
	if err != nil {
		panic(err)
	}
	if err = c.AddAuth("digest", []byte(user+":"+password)); err != nil {
		log.Printf("AddAuth user returned error\n")
	}
	if openUser != "" {
		if err = c.AddAuth("digest", []byte(openUser+":"+openPassword)); err != nil {
			log.Printf("AddAuth openUser returned error")
		}
	}
	cl = &zkClient{c: c}
	return
}

func (c *zkClient) GetValues(keys []string) (vars map[string]string, err error) {
	vars = make(map[string]string)
	for _, v := range keys {
		_, _, err = c.c.Exists(v)
		if err != nil {
			return
		}
		b, _, err := c.c.Get(v)
		if err != nil {
			return vars, err
		}
		vars[v] = string(b)
	}
	return
}

func (c *zkClient) WatchPrefix(keys []string, waitIndex uint64) (index uint64, err error) {
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

func (c *zkClient) watch(key string, respChan chan WatchResponse, cancelRoutine chan bool) {
	_, _, keyEventCh, err := c.c.GetW(key)
	if err != nil {
		respChan <- WatchResponse{Err: err}
	}
	_, _, childEventCh, err := c.c.ChildrenW(key)
	if err != nil {
		respChan <- WatchResponse{Err: err}
	}

	for {
		select {
		case e := <-keyEventCh:
			if e.Type == zk.EventNodeDataChanged {
				respChan <- WatchResponse{WaitIndex: 1, Err: e.Err}
			}
		case e := <-childEventCh:
			if e.Type == zk.EventNodeChildrenChanged {
				respChan <- WatchResponse{WaitIndex: 1, Err: e.Err}
			}
		case <-cancelRoutine:
			log.Printf("Stop watching:%v\n", key)
			return
		}
	}
}
