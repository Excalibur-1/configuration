package configuration

import (
	"fmt"
	"testing"
)

func TestDefaultEngine(t *testing.T) {
	engine := DefaultEngine()
	s, err := engine.String("myconf", "base", "cache", "", "1000")
	fmt.Println(s, err)
}

func TestMockEngine(t *testing.T) {
	engine := MockEngine(map[string]string{
		"/myconf/base/cache/1000": "{\"provider\":\"redis\"}",
	})
	s, err := engine.String("myconf", "base", "cache", "", "1000")
	fmt.Println(s, err)
}

func TestZkEngine(t *testing.T) {
	engine := ZkEngine(NewStoreConfig(""))
	s, err := engine.String("myconf", "base", "cache", "", "1000")
	fmt.Println(s, err)
}

func TestEtcdEngine(t *testing.T) {
	engine := EtcdEngine(NewStoreConfig(""))
	s, err := engine.String("myconf", "base", "cache", "", "1000")
	fmt.Println(s, err)
}
