// Package configuration 提供给所有业务系统使用的配置管理引擎，所有业务系统/中间件的可变配置通过集中配置中心进行统一配置，
// 业务系统可通过该引擎来管理自己的配置数据并在配置中心的数据发生变化时得到及时的通知。
// 配置管理中心的数据存储为树形目录结构（类似文件系统的目录结构），每一个节点都可以存储相应的数据。业务系统可在基础目录树的基础上
// 往下继续扩展新的目录结构，从而可以区分存储不同的配置项信息。比如系统分配给业务系统A的基础目录（称为path）为：/sysa，
// 那么业务系统可以扩充到如下path：/sysa/key1、/sysa/key2、/sysa/key3，每个path可存储对应的数据。
// 配置管理引擎客户端依赖连接的配置管理中心，需要在统一配置管理中心下载对应的授权文件方可。
package configuration

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type Configuration interface {
	Values(namespace, app, group, tag string, path []string) (map[string]string, error)
	String(namespace, app, group, tag, path string) (string, error)
	Clazz(namespace, app, group, tag, path string, clazz interface{}) error
	Get(namespace, app, group, tag string, path []string, parser ChangedListener)
}

type configuration struct {
	store StoreClient
}

// DefaultEngine 默认采用zk配置中心且uaf环境变量设置为：UAF
func DefaultEngine() Configuration {
	return ZkEngine(NewStoreConfig("", ""))
}

// MockEngine 本地mock一个配置
func MockEngine(conf map[string]string) Configuration {
	fmt.Println("Loading configuration MockEngine ver:1.0.0")
	store, err := NewMockClient(conf)
	if err != nil {
		panic(err)
	}
	return &configuration{store: store}
}

// ZkEngine 获取zk配置管理引擎的唯一实例。
func ZkEngine(conf StoreConfig) Configuration {
	fmt.Println("Loading configuration ZkEngine ver:1.0.0")
	store, err := NewZkClient(conf.Servers, conf.Username, conf.Password, conf.OpenUser, conf.OpenPassword)
	if err != nil {
		panic(err)
	}
	return &configuration{store: store}
}

// EtcdEngine 获取etcd配置管理引擎的唯一实例。
func EtcdEngine(conf StoreConfig) Configuration {
	fmt.Println("Loading configuration EtcdEngine ver:1.0.0")
	store, err := NewEtcdClient(conf.Servers, conf.Username, conf.Password)
	if err != nil {
		panic(err)
	}
	return &configuration{store: store}
}

// Values 获取多个配置项的配置信息，返回原始的配置数据格式(map集合)，如果获取失败则抛出异常。
func (c configuration) Values(namespace, app, group, tag string, path []string) (map[string]string, error) {
	_path := make([]string, len(path))
	for i, v := range path {
		_path[i] = c.maskPath(namespace, app, group, tag, v)
	}
	vl, err := c.store.GetValues(_path)
	if err != nil {
		log.Printf("获取多个配置项[%v]的配置信息出错\n", path)
	} else {
		log.Printf("获取多个配置项为:%+v\n", vl)
	}
	return vl, err
}

// String 获取指定配置项的配置信息，返回原始的配置数据格式，如果获取失败则抛出异常。
func (c configuration) String(namespace, app, group, tag, path string) (string, error) {
	path = c.maskPath(namespace, app, group, tag, path)
	vl, err := c.store.GetValues([]string{path})
	if err != nil {
		log.Printf("获取多个配置项[%s]的配置信息出错\n", path)
	} else {
		log.Printf("获取多个配置项为:%+v\n", vl)
	}
	return vl[path], err
}

// Clazz 获取指定配置项的配置信息，并且将配置信息（JSON格式的）转换为指定的Go结构体，如果获取失败或转换失败则抛出异常。
func (c configuration) Clazz(namespace, app, group, tag, path string, clazz interface{}) error {
	path = c.maskPath(namespace, app, group, tag, path)
	vl, err := c.store.GetValues([]string{path})
	if err != nil {
		log.Printf("获取多个配置项[%s]的配置信息出错:%+v\n", path, err)
		return err
	} else {
		log.Printf("获取配置项为:%+v\n", vl)
	}
	err = json.Unmarshal([]byte(vl[path]), clazz)
	return err
}

// Get 获取指定路径下的配置信息，并实现监听，当有数据变化时自动调用解析器进行解析。
func (c configuration) Get(namespace, app, group, tag string, path []string, parser ChangedListener) {
	_path := make([]string, len(path))
	for i, v := range path {
		_path[i] = c.maskPath(namespace, app, group, tag, v)
	}
	vl, err := c.store.GetValues(_path)
	if err != nil {
		log.Printf("获取指定路径[%v]下的配置信息,并实现监听,当有数据变化时自动调用解析器进行解析出错\n", path)
	} else {
		log.Printf("获取多个配置项为:%+v\n", vl)
	}
	parser.Changed(vl)
	NewProcessor(_path, c.store).Process(parser)
}

func (configuration) maskPath(namespace, app, group, tag, path string) string {
	key := []string{"/" + namespace, app, group, path}
	if len(tag) > 0 {
		key = []string{"/" + namespace, app, group, tag, path}
	}
	return strings.Join(key, "/")
}
