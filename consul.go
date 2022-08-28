package aurora

import (
	"errors"
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
func (engine *Engine) consul() {

	// 对对基本配置进行检查，确保后续正确性
	if preCheck, err := engine.preCheck(); !preCheck {
		ErrorMsg(err)
		return
	}
	consulConfigs := engine.config.GetStringMapString("aurora.consul")
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
	ErrorMsg(err)
	// 解析注册地址
	hosts := strings.Split(registers, ",")
	// 创建 每个 consul 的 客户端
	consuls := engine.consulHost(hosts)
	// 创建失败 则结束配置
	if consuls == nil {
		return
	}
	engine.consulCenter = &ConsulCenter{consuls: consuls}
	// conf 若配置 则添加远程配置地址 并覆盖本地配置环境
	if conf != "" {
		v := viper.New()
		v.SetConfigType("yaml")
		for host, consul := range consuls {
			consul.Address = host
			consul.Config = conf
			err = v.AddRemoteProvider("consul", host, conf)
			ErrorMsg(err)
		}
		err = v.ReadRemoteConfig()
		if err != nil {
			ErrorMsg(err)
		}
		// 刷新本地配置
		cnf := &ConfigCenter{
			v,
			&sync.RWMutex{},
		}
		engine.config = cnf
		// 配置文件 监听
		go func(center *ConfigCenter) {
			for true {
				time.Sleep(time.Duration(refresh) * time.Second)
				//old := center.GetStringMap("aurora")
				err = center.WatchRemoteConfig()
				if err != nil {
					engine.Error(err.Error())
					continue
				}
				//new := center.GetStringMap("aurora")
			}
		}(cnf)
	}

	// 生成 web 服务
	registration := engine.getAgentServiceRegistration()
	// 向客户端注册服务
	for _, consul := range consuls {
		err = consul.Agent.ServiceRegister(registration)
		if err != nil {
			engine.Error(err.Error())
		}

	}

	// consul 配置完毕 把 consul 的配置中心加入到 ioc 中
	engine.Use(engine.consulCenter)
}

// 创建集群客户端
func (engine *Engine) consulHost(hosts []string) map[string]*Consul {
	consuls := make(map[string]*Consul)
	config := engine.getConsulClientConfig()
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

// 初始化 consul客户端公共配置
func (engine *Engine) getConsulClientConfig() *api.Config {
	config := api.DefaultConfig()
	// 读取 客户端初始化配置项
	clientConfig := engine.config.GetStringMapString("aurora.consul.client")
	for key, value := range clientConfig {
		if value != "" {
			c := keymap[key]
			configMap[c](value)(config)
		}
	}
	return config
}

// 生成当前 web 服务注册信息
func (engine *Engine) getAgentServiceRegistration() *api.AgentServiceRegistration {
	// 读取 服务 名称
	name := engine.config.GetString("aurora.server.name")

	// 读取 服务 ip地址
	host := engine.config.GetString("aurora.server.host")

	// 读取 服务 端口
	port := engine.config.GetString("aurora.server.port")

	atoi, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}

	// 生成唯一服务id 区分集群
	format := fmt.Sprintf("%s-%s:%s", name, host, port)

	// 创建服务
	registration := &api.AgentServiceRegistration{
		ID:      format,
		Name:    name,
		Port:    atoi,
		Address: host,
		Check:   engine.getAgentServiceCheck(),
	}
	return registration
}

func (engine *Engine) getAgentServiceCheck() *api.AgentServiceCheck {

	// 读取服务检查 地址 aurora 默认采用 http 方式
	url := engine.config.GetString("aurora.consul.service.check.url")

	// 读取 心跳检查频率
	interval := engine.config.GetString("aurora.consul.service.check.interval")

	// 读取 服务超时时间
	timeout := engine.config.GetString("aurora.consul.service.check.timeout")

	checkName := engine.config.GetString("aurora.server.name")

	// 读取检查名称
	if name := engine.config.GetString("aurora.consul.service.check.name"); name == "" {
		//生成服务检查名称
		checkName = fmt.Sprintf("Service '%s' check", checkName)
	} else {
		checkName = name
	}

	// 读取 服务 名称
	name := engine.config.GetString("aurora.server.name")

	// 读取 服务 ip地址
	host := engine.config.GetString("aurora.server.host")

	// 读取 服务 端口
	port := engine.config.GetString("aurora.server.port")

	// 生成检查ID
	checkId := fmt.Sprintf("Service:%s-%s:%s", strings.ToUpper(name), host, port)

	// 创建 服务检查
	c := &api.AgentServiceCheck{
		CheckID:       checkId,
		Name:          checkName,
		Interval:      interval,
		Timeout:       timeout,
		HTTP:          url,
		Method:        http.MethodGet, // 默认使用post 进行检查
		TLSSkipVerify: true,           // 默认不开启 tls 检查
	}
	return c
}

// 配置consul 之前的配置预检查
func (engine *Engine) preCheck() (bool, error) {
	// 检查服务名是否配置

	s := engine.config.Get("aurora.consul")
	if s == nil {
		return false, nil
	}
	// 检查 是否启用 consul
	getString := engine.config.GetString("aurora.consul.enable")

	// 没有配置 enable 默认启动
	if getString != "" {
		enable, err := strconv.ParseBool(getString)
		if err != nil {
			return false, err
		}

		// 不启动 返回 false 和 nil 不进行下面的检查
		if !enable {
			return enable, nil
		}
	}
	// 读取 服务 名称
	if name := engine.config.GetString("aurora.server.name"); name == "" {
		return false, errors.New("no service name is configured, please check the configuration file configuration item 'aurora.server.name'")
	}

	// 检查 端口号
	// 读取 服务 端口
	if port := engine.config.GetString("aurora.server.port"); port == "" {
		return false, errors.New("no service port is configured, please check the configuration file configuration item 'aurora.server.port'")
	}

	// 检查 主机号
	// 读取 服务 ip地址
	if host := engine.config.GetString("aurora.server.host"); host == "" {
		return false, errors.New("no service host is configured, please check the configuration file configuration item 'aurora.server.host'")
	}

	return true, nil
}

// Health consul 健康检查回掉
func Health() string {
	return "ok"
}

func (c *Consul) GetService() {

}
