package aurora

import (
	"github.com/hashicorp/consul/api"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Consul struct {
	Host   string // 当前consul 远程地址
	Config string // consul KV 读取中心
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

func (a *Aurora) consul() {
	consulConfig := a.config.GetStringMapString("aurora.consul")
	if consulConfig == nil {
		return
	}
	registers := consulConfig["register"]
	// 解析 地址
	hosts := strings.Split(registers, ",")
	consuls := consulHost(hosts)
	if consuls == nil {
		return
	}
	// 生成 web 服务
	registration := a.getAgentServiceRegistration()
	for _, consul := range consuls {
		err := consul.Agent.ServiceRegister(registration)
		if err != nil {
			a.Error(err.Error())
		}
	}
}

func consulHost(hosts []string) map[string]*Consul {
	consuls := make(map[string]*Consul)
	config := api.DefaultConfig()
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

// 生成当前 web 服务注册信息
func (a *Aurora) getAgentServiceRegistration() *api.AgentServiceRegistration {
	name := a.config.GetString("aurora.server.name")
	host := a.config.GetString("aurora.server.host")
	port := a.config.GetString("aurora.server.port")
	atoi, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	registration := &api.AgentServiceRegistration{
		ID:      name,
		Name:    name,
		Port:    atoi,
		Address: host,
		Check:   a.getAgentServiceCheck(),
	}
	return registration
}

func (a *Aurora) getAgentServiceCheck() *api.AgentServiceCheck {
	url := a.config.GetString("aurora.consul.service.check.url")
	name := a.config.GetString("aurora.consul.service.check.name")
	interval := a.config.GetString("aurora.consul.service.check.interval")
	timeout := a.config.GetString("aurora.consul.service.check.timeout")
	c := &api.AgentServiceCheck{
		CheckID:       "",
		Name:          name,
		Interval:      interval,
		Timeout:       timeout,
		HTTP:          url,
		Method:        http.MethodGet,
		TLSSkipVerify: true,
	}
	return c
}
