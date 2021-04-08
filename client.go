package configuration

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Excalibur-1/gutil"
	"io/ioutil"
	"os"
)

const (
	DesKey      = "DES_KEY"
	DesKeyValue = "12345678"
	Uaf         = "UAF"
	UafFileName = "./conf/configuration.uaf"
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

// NewStoreConfig 读取配置中心配置，如果没有传递环境变量，则取默认值
func NewStoreConfig(uafEnv, desKeyEnv string) (conf StoreConfig) {
	var uaf, desKey string
	if uafEnv == "" {
		uaf = os.Getenv(Uaf)
	} else {
		uaf = os.Getenv(uafEnv)
	}
	if desKeyEnv == "" {
		desKey = os.Getenv(DesKey)
	} else {
		desKey = os.Getenv(desKeyEnv)
	}
	if uaf == "" {
		b, err := ioutil.ReadFile(UafFileName)
		if err != nil {
			panic("未找到配置管理授权文件：" + UafFileName)
		}
		uaf = string(b)
	}
	if desKey == "" {
		desKey = DesKeyValue
	}
	if uaf == "" {
		panic(fmt.Sprintf("请在环境变量中配置UAF或者%s中配置内容", UafFileName))
	}
	ds, err := base64.URLEncoding.DecodeString(uaf)
	if err != nil {
		panic("解码configuration.uaf出错")
	}
	data, err := gutil.Decrypt(ds, []byte(desKey))
	if err != nil {
		panic("解密configuration.uaf出错")
	}
	if err = json.Unmarshal(data, &conf); err != nil {
		panic("不存在配置资源请使用本地配置信息")
	}
	return
}
