package configuration

import (
	"encoding/base64"
	"encoding/json"
	"github.com/Excalibur-1/gutil"
	"io/ioutil"
	"os"
)

const (
	DesKey = "12345678"
)

// The StoreClient interface is implemented by objects that can retrieve key/value pairs from a backend store.
type StoreClient interface {
	GetValues(keys []string) (map[string]string, error)
	WatchPrefix(keys []string, waitIndex uint64) (uint64, error)
}

type StoreConfig struct {
	Servers      []string          `json:"servers"`      // 配置中心节点地址列表
	Username     string            `json:"username"`     // 只读用户名
	Password     string            `json:"password"`     // 只读用户密码
	OpenUser     string            `json:"openUser"`     // 可读写用户名
	OpenPassword string            `json:"openPassword"` // 可读写密码
	Exp          map[string]string `json:"exp"`
}

type WatchResponse struct {
	WaitIndex uint64
	Err       error
}

func NewStoreConfig() (conf StoreConfig) {
	uaf := os.Getenv("UAF")
	if uaf == "" {
		b, err := ioutil.ReadFile("./conf/configuration.uaf")
		if err != nil {
			panic("未找到配置管理授权文件：./conf/configuration.uaf")
		}
		uaf = string(b)
	}
	if uaf == "" {
		panic("请在环境变量中配置UAF或者./conf/configuration.uaf中配置内容")
	}
	ds, err := base64.URLEncoding.DecodeString(uaf)
	if err != nil {
		panic("解码configuration.uaf出错")
	}
	data, err := gutil.Decrypt(ds, []byte(DesKey))
	if err != nil {
		panic("解密configuration.uaf出错")
	}
	if err = json.Unmarshal(data, &conf); err != nil {
		panic("不存在配置资源请使用本地配置信息")
	}
	return
}
