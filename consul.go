package aurora

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Consul struct {
	Address string // 当前consul 远程地址
	Config  string // consul KV 读取中心
	*api.Client
	*api.Agent
	*api.KV
}

func newConsul(client *api.Client) *Consul {
	c := &Consul{}
	c.Client = client
	c.Agent = client.Agent()
	c.KV = client.KV()
	return c
}

// 读取 配置文件 aurora.consul 并配置
func (a *Aurora) consul() {
	getString := a.config.GetString("enable")
	if getString != "" {
		enable, err := strconv.ParseBool(getString)
		if err != nil {
			log.Print(err.Error())
			return
		}
		if !enable {
			return
		}
	}

	consulConfigs := a.config.GetStringMapString("aurora.consul")
	if consulConfigs == nil {
		return
	}

	// 找到 consul 服务地址，用于注册本服务
	registers := consulConfigs["register"]

	// config 读取 consul k/v 用于指定的配置文件读取
	// 若是没有配置，则继续使用本地配置
	conf := consulConfigs["config"]
	ref := consulConfigs["refresh"]
	if ref[len(ref)-1:] != "s" {
		panic("configuration refresh time format error")
	}
	refresh, err := strconv.Atoi(ref[:len(ref)-1])
	if err != nil {
		panic(err)
		return
	}
	// 解析注册地址
	hosts := strings.Split(registers, ",")
	// 创建 每个 consul 的 客户端
	consuls := a.consulHost(hosts)

	// 创建失败 则结束配置
	if consuls == nil {
		return
	}

	// conf 若配置 则添加远程配置地址 并覆盖本地配置环境
	if conf != "" {
		v := viper.New()
		v.SetConfigType("yaml")
		for host, consul := range consuls {
			consul.Address = host
			consul.Config = conf
			err = v.AddRemoteProvider("consul", host, conf)
			if err != nil {
				panic(err)
				return
			}
		}
		err = v.ReadRemoteConfig()
		if err != nil {
			panic(err)
			return
		}
		// 刷新本地配置
		cnf := &ConfigCenter{
			v,
			&sync.RWMutex{},
		}
		a.config = cnf
		// 配置文件 监听
		go func(center *ConfigCenter) {
			for true {
				time.Sleep(time.Duration(refresh) * time.Second)
				//old := center.GetStringMap("aurora")
				err = center.WatchRemoteConfig()
				if err != nil {
					a.Error(err.Error())
					continue
				}
				//new := center.GetStringMap("aurora")
			}
		}(cnf)
	}

	// 生成 web 服务
	registration := a.getAgentServiceRegistration()
	// 向客户端注册服务
	for _, consul := range consuls {
		err = consul.Agent.ServiceRegister(registration)
		if err != nil {
			a.Error(err.Error())
		}

	}

}

// 创建集群客户端
func (a *Aurora) consulHost(hosts []string) map[string]*Consul {
	consuls := make(map[string]*Consul)
	config := a.getConsulClientConfig()
	for _, host := range hosts {
		if host != "" {
			config.Address = host
			client, err := api.NewClient(config)
			if err != nil {
				log.Fatal(err)
				return nil
			}
			consuls[host] = newConsul(client)
		}
	}
	return consuls
}

func (a *Aurora) getConsulClientConfig() *api.Config {
	config := api.DefaultConfig()
	// 读取 客户端初始化配置项
	clientConfig := a.config.GetStringMapString("aurora.consul.client")
	for key, value := range clientConfig {
		if value != "" {
			c := keymap[key]
			configMap[c](value)(config)
		}
	}
	return config
}

// 生成当前 web 服务注册信息
func (a *Aurora) getAgentServiceRegistration() *api.AgentServiceRegistration {
	// 读取 服务 名称
	name := a.config.GetString("aurora.server.name")

	// 读取 服务 ip地址
	host := a.config.GetString("aurora.server.host")

	// 读取 服务 端口
	port := a.config.GetString("aurora.server.port")

	atoi, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}

	// 生成唯一服务id
	format := fmt.Sprintf("%s-%s:%s", strings.ToUpper(name), host, port)

	// 创建服务
	registration := &api.AgentServiceRegistration{
		ID:      format,
		Name:    name,
		Port:    atoi,
		Address: host,
		Check:   a.getAgentServiceCheck(),
	}
	return registration
}

func (a *Aurora) getAgentServiceCheck() *api.AgentServiceCheck {

	// 读取服务检查 地址 aurora 默认采用 http 方式
	url := a.config.GetString("aurora.consul.service.check.url")

	// 读取服务名称
	name := a.config.GetString("aurora.consul.service.check.name")

	// 读取 心跳检查频率
	interval := a.config.GetString("aurora.consul.service.check.interval")

	// 读取 服务超时时间
	timeout := a.config.GetString("aurora.consul.service.check.timeout")

	// 创建 服务检查
	c := &api.AgentServiceCheck{
		CheckID:       "",
		Name:          name,
		Interval:      interval,
		Timeout:       timeout,
		HTTP:          url,
		Method:        http.MethodGet, // 默认使用post 进行检查
		TLSSkipVerify: true,           // 默认不开启 tls 检查
	}
	return c
}

// Health consul 健康检查回掉
func Health() string {
	return "ok"
}
