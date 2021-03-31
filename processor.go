package configuration

import (
	"log"
	"time"
)

// ChangedListener 配置数据发生变化时的监听器接口。
type ChangedListener interface {
	// Changed 指定配置变化后的通知接口，标示变化的路径和变化后的数据。
	Changed(data map[string]string)
}

type Processor interface {
	Process(listener ChangedListener)
}

type watchProcessor struct {
	path  []string
	store StoreClient
}

func NewProcessor(path []string, store StoreClient) Processor {
	return &watchProcessor{
		path:  path,
		store: store,
	}
}

func (w *watchProcessor) Process(listener ChangedListener) {
	go w.monitorPrefix(w.path, listener)
}

func (w *watchProcessor) monitorPrefix(path []string, listener ChangedListener) {
	var lastIndex uint64
	for {
		index, err := w.store.WatchPrefix(path, lastIndex)
		if err != nil {
			log.Printf("monitorPrefix WatchPrefix has err:%v\n", err)
			time.Sleep(time.Second * 2)
			continue
		}
		if lastIndex > 0 {
			if vl, er := w.store.GetValues(path); er != nil {
				log.Printf("monitorPrefix GetValues has err:%v\n", err)
			} else {
				listener.Changed(vl)
			}
		}
		lastIndex = index
	}
}
