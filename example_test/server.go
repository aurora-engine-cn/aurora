package example

import (
	"fmt"
	"gitee.com/aurora-engine/aurora"
)

// Server 嵌套Aurora定义一个服务 实例
type Server struct {
	*aurora.Engine
}

func (server *Server) Server() {
	// 进行一下初始化操作，比如 控制器实例，全局中间件，全局变量，第三方依赖库等操作
}

func (server *Server) Router() {
	// 添加 app 路由

	server.Get("/a/{name}/bbb/{age}/update", func(name, age string) string {
		fmt.Printf("A name:%s,age:%s\n", name, age)
		return "hello world"
	})

	server.Get("/b/{name}/bbb/{age}/update", func(name, age string) string {
		fmt.Printf("B name:%s,age:%s\n", name, age)
		return "hello world"
	})
	server.Get("/c/{name}/{age}/update", func(name, age string) string {
		fmt.Printf("C1 name:%s,age:%s\n", name, age)
		return "hello world"
	})

	//server.Get("/c/{name}/{age}/update", func(name, age string) string {
	//	fmt.Printf("C2 name:%s,age:%s\n", name, age)
	//	return "hello world"
	//})

}
