package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Excalibur-1/configuration"
	"github.com/Excalibur-1/gutil"
	"io/ioutil"
	"sync"
)

// 测试的配置文件为格式为：
// "{\"servers\":[\"localhost:2181\"],\"username\":\"guest\",\"password\":\"guest\"}"
func main() {
	go zk()
	go etcd()
	makeUaf()
}

func zk() {
	zkChange()
	for {
		s, err := configuration.DefaultEngine().String("myconf", "base", "cache", "", "provider")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(s)
		var str string
		_, _ = fmt.Scanln(&str)
		fmt.Println(providerCfg)
	}
}

func etcd() {
	etcdChange()
	for {
		s, err := configuration.
			EtcdEngine(configuration.NewStoreConfig("")).
			String("myconf", "base", "cache", "", "provider")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(s)
		var str string
		_, _ = fmt.Scanln(&str)
		fmt.Println(providerCfg)
	}
}

type ProviderConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

var providerCfg *ProviderConfig

func zkChange() {
	providerCfg = &ProviderConfig{}
	eng := configuration.DefaultEngine()
	if s, err := eng.String("myconf", "base", "cache", "", "provider"); err != nil {
		_ = json.Unmarshal([]byte(s), providerCfg)
	}
	eng.Get("myconf", "base", "cache", "", []string{"provider"}, &providerParser{})
}

func etcdChange() {
	providerCfg = &ProviderConfig{}
	eng := configuration.EtcdEngine(configuration.NewStoreConfig(""))
	if s, err := eng.String("myconf", "base", "cache", "", "provider"); err != nil {
		_ = json.Unmarshal([]byte(s), providerCfg)
	}
	eng.Get("myconf", "base", "cache", "", []string{"provider"}, &providerParser{})
}

type providerParser struct {
	mu sync.Mutex
}

func (t *providerParser) Changed(data map[string]string) {
	fmt.Printf("Changed %+v\n", data)
	for _, v := range data {
		var vl ProviderConfig
		if err := json.Unmarshal([]byte(v), &vl); err != nil {
			fmt.Printf("new value for config error: %+v\n", err)
			return
		}
		t.mu.Lock()
		providerCfg = &vl
		fmt.Printf("providerCfg changed to: %+v\n", providerCfg)
		t.mu.Unlock()
	}
}

func makeUaf() {
	key := "12345678"
	src := "{\"servers\":[\"localhost:2181\"],\"username\":\"guest\",\"password\":\"guest\"}"
	fmt.Println("原文：" + src)
	enc, _ := gutil.Encrypt([]byte(src), []byte(key))
	_ = ioutil.WriteFile("configuration.uaf", []byte(base64.URLEncoding.EncodeToString(enc)), 0644)
	fmt.Println("密文：" + base64.URLEncoding.EncodeToString(enc))
	decodeString, _ := base64.URLEncoding.DecodeString("vs13PY3sDZ6ryktf5akSn6GBIJF-oC4II-Pm5x0wxY5mTQVfwCmTddt4ugBr1nwSBmEmlwjCLlS1jrcL9FbWRUh-LjWhtJgpUZLbo-akf3A=")
	dec, _ := gutil.Decrypt(decodeString, []byte(key))
	fmt.Println("解码：" + string(dec))
}
