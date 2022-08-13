package aurora

import (
	"github.com/spf13/viper"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
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

// viperConfig 配置并加载 application.yml 配置文件
func (a *Aurora) viperConfig() {
	var ConfigPath string

	// consul 配置中心校验

	if a.configpath == "" {
		//检索配置文件所在路径
		err := filepath.WalkDir(a.projectRoot, func(p string, d fs.DirEntry, err error) error {
			//找到配置及文件,基于根目录优先加载最外层的application.yml
			if !d.IsDir() && (strings.HasSuffix(p, yml) || (strings.HasSuffix(p, yaml))) {
				//修复 项目加载配置覆盖，检索项目配置文件，避免内层同名配置文件覆盖外层，这个情况可能发生在 开发者把两个go mod 项目嵌套在一起，导致配置被覆盖
				//此处校验，根据检索的更路径，只加载最外层的配置文件
				if ConfigPath == "" {
					ConfigPath = p
				}
			}
			return nil
		})
		ErrorMsg(err, "configuration file retrieval failed, initial configuration failed")
	} else {
		ConfigPath = a.configpath
	}
	if ConfigPath == "" {
		a.config = &ConfigCenter{viper.New(), &sync.RWMutex{}}
		return
	}
	if a.config == nil {
		// 用户没有提供 配置项 则创建默认的配置处理
		cnf := &ConfigCenter{
			viper.New(),
			&sync.RWMutex{},
		}
		cnf.SetConfigFile(ConfigPath)
		err := cnf.ReadInConfig()
		ErrorMsg(err)
		a.config = cnf
	}

}

// GetConfig 获取 Aurora 配置实例 对配置文件内容的读取都是协程安全的
func (a *Aurora) GetConfig() Config {
	return a.config
}

// ViperConfiguration 配置指定配置文件
func ViperConfiguration(path string) Option {
	return func(a *Aurora) {

	}
}
