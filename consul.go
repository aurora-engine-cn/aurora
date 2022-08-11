package aurora

import "github.com/hashicorp/consul/api"

type Consul struct {
	Host   string // 当前consul 远程地址
	Config string // consul KV 读取中心
	*api.Client
	*api.Agent
	*api.KV
}

func (c *Consul) config() {

}
