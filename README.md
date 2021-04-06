# configuration
配置中心

## 快速使用
1. 生成 uaf 文件
    * 根据需要修改 example 目录下的 makeUaf() 方法，修改配置中心服务器地址及账号密码，然后在main方法中调用即可； 
    * 上一步执行完成将得到 configuration.uaf 文件，将其添加到项目所在目录的 conf 文件夹下；
    * 或者把 configuration.uaf 文件中的内容设置到环境变量中，环境变量的key为指定的值，默认是：UAF
    
2. 调用实例(默认配置，用的是zookeeper)
```go
import (
	"fmt"
	"github.com/Excalibur-1/configuration"
)

func example() {
	engine := configuration.DefaultEngine()
	value, err := engine.String("myconf", "base", "cache", "", "1000")
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
}

```
上面采用默认的配置中心zk来作为配置中心使用，且 uaf内容 读取的环境变量值也是默认的：`UAF`

3. 自定义配置
```go
import (
	"fmt"
	"github.com/Excalibur-1/configuration"
)

func example() {
	storeCfg := configuration.NewStoreConfig("MY_UAF")
	engine := configuration.EtcdEngine(storeCfg)
	value, err := engine.String("myconf", "base", "cache", "", "1000")
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
}

```
上面是采用自定义的配置中心来读取数据，首先传入 uaf 环境变量的key，然后将读取到的配置传递给 EtcdEngine 方法获取实例

4. 监听配置变更
```go
import (
	"encoding/json"
	"fmt"
	"github.com/Excalibur-1/configuration"
	"sync"
)

func init() {
	providerCfg = &ProviderConfig{}
	eng := configuration.DefaultEngine()
	if s, err := eng.String("myconf", "base", "cache", "", "1000"); err != nil {
		_ = json.Unmarshal([]byte(s), providerCfg)
	}
	eng.Get("myconf", "base", "cache", "", []string{"1000"}, &providerParser{})
}

var providerCfg *ProviderConfig

type ProviderConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
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

```
只要实现 Change 方法，就可定义配置变更时需要做的处理，这里就是直接对全局变量的值进行修改。

## 2021-03-31 更新日志
* 完成zookeeper和etcd配置中心服务封装

## 2021-04-06 更新日志
* 代码结构调整及优化
* 添加 readme 使用说明文档
* 添加单元测试用例