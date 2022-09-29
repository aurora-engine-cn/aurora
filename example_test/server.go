package example

import "gitee.com/aurora-engine/aurora"

// Server 嵌套Aurora定义一个服务 实例
type Server struct {
	*aurora.Engine
}

func (server *Server) Server() {
	// 进行一下初始化操作，比如 控制器实例，全局中间件，全局变量，第三方依赖库等操作
}

func (server *Server) Router() {
	// 添加 app 路由

	server.Get("/", func() string {
		return "hello world"
	})
}
