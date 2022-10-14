package utils

//import (
//	"errors"
//	"fmt"
//	"github.com/hashicorp/consul/api"
//)
//
//// ConsulCenter consul 管理中心
//type ConsulCenter struct {
//	consuls map[string]*Consul
//}
//
//// Services 获取 consul 中的所有服务信息
//func (c *ConsulCenter) Services() (map[string]*api.AgentService, error) {
//	return c.ServicesWithFilter("")
//}
//
//func (c *ConsulCenter) ServicesWithFilter(filter string) (map[string]*api.AgentService, error) {
//	return c.ServicesWithFilterOpts(filter, nil)
//}
//
//func (c *ConsulCenter) ServicesWithFilterOpts(filter string, q *api.QueryOptions) (map[string]*api.AgentService, error) {
//	for _, consul := range c.consuls {
//		services, err := consul.ServicesWithFilterOpts(filter, q)
//		if err != nil {
//			fmt.Println(err)
//			continue
//		} else {
//			return services, err
//		}
//	}
//	return nil, errors.New("no services found")
//}
