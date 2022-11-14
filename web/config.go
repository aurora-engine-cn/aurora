package web

import (
	"github.com/spf13/viper"
	"io"
	"sync"
)

const (
	yml  = "application.yml"
	yaml = "application.yaml"
)

/*
	viper 配置文件实例
	提供Aurora 默认的配置实例，
	默认读取配置文件的位置为根目录 application.yml application.yaml
*/

type Config interface {
	SetConfigFile(string)
	SetConfigType(string)
	ReadConfig(io.Reader) error
	Set(string, interface{})
	SetDefault(string, interface{})
	GetStringMapString(string) map[string]string
	Get(string) interface{}
	GetStringSlice(string) []string
	GetStringMap(string) map[string]interface{}
	GetString(string) string
	GetStringMapStringSlice(string) map[string][]string
}

// ConfigCenter 配置中心 的读写锁主要用来解决分布式配置的动态刷新配置，和以后存在的并发读取配置和修改，
// 对于修改配置数据库连接信息或者需要重新初始化的配置项这些无法起到同步更新的效果只能保持配置信息是最新的（需要重新初始化的配置建议重启服务），
// 对被配置的使用实例没有并发安全的效果。
type ConfigCenter struct {
	*viper.Viper
	*sync.RWMutex
}

func (c *ConfigCenter) SetConfigFile(in string) {
	c.Lock()
	defer c.Unlock()
	c.Viper.SetConfigFile(in)
}

func (c *ConfigCenter) SetConfigType(in string) {
	c.Lock()
	defer c.Unlock()
	c.Viper.SetConfigType(in)
}
func (c *ConfigCenter) ReadConfig(in io.Reader) error {
	c.Lock()
	defer c.Unlock()
	return c.Viper.ReadConfig(in)
}

func (c *ConfigCenter) Set(key string, value interface{}) {
	c.Lock()
	defer c.Unlock()
	c.Viper.Set(key, value)
}

func (c *ConfigCenter) SetDefault(key string, value interface{}) {
	c.Lock()
	defer c.Unlock()
	c.Viper.SetDefault(key, value)
}

//读取配置文件

func (c *ConfigCenter) ReadInConfig() error {
	c.RLock()
	defer c.RUnlock()
	return c.Viper.ReadInConfig()
}

func (c *ConfigCenter) WatchRemoteConfig() error {
	c.Lock()
	defer c.Unlock()
	return c.Viper.WatchRemoteConfig()
}

func (c *ConfigCenter) GetStringMapString(key string) map[string]string {
	c.RLock()
	defer c.RUnlock()
	return c.Viper.GetStringMapString(key)
}

func (c *ConfigCenter) Get(key string) interface{} {
	c.RLock()
	defer c.RUnlock()
	return c.Viper.Get(key)
}

func (c *ConfigCenter) GetStringSlice(key string) []string {
	c.RLock()
	defer c.RUnlock()
	return c.Viper.GetStringSlice(key)
}

func (c *ConfigCenter) GetStringMap(key string) map[string]interface{} {
	c.RLock()
	defer c.RUnlock()
	return c.Viper.GetStringMap(key)
}

func (c *ConfigCenter) GetString(key string) string {
	c.RLock()
	defer c.RUnlock()
	return c.Viper.GetString(key)
}

func (c *ConfigCenter) GetStringMapStringSlice(key string) map[string][]string {
	c.RLock()
	defer c.RUnlock()
	return c.Viper.GetStringMapStringSlice(key)
}
