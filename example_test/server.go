package example

import (
	"errors"
	"gitee.com/aurora-engine/aurora"
	"net/http/pprof"
)

// Server 嵌套Aurora定义一个服务 实例
type Server struct {
	*aurora.Engine
}

type GetArgs struct {
	Name string `empty:"false"`
	Age  int    `constraint:"check"`
}

func (server *Server) Server() {
	// 进行一下初始化操作，比如 控制器实例，全局中间件，全局变量，第三方依赖库等操作

	server.Constraint("check", func(value any) error {
		if value.(int) <= 0 {
			return errors.New("error value is 0")
		}
		return nil
	})
}

func (server *Server) Router() {
	// 添加 app 路由

	server.Get("test", func(args GetArgs) {

	})

	server.Post("/user", func(name, age string) string {
		return ""
	})
	server.Get("/user/{id}", func(id string) string {
		return id
	})

	pprofs := server.Group("/debug")
	pprofs.Get("/pprof", pprof.Index)
	pprofs.Get("/pprof/cmdline", pprof.Cmdline)
	pprofs.Get("/pprof/profile", pprof.Profile)
	pprofs.Get("/pprof/symbol", pprof.Symbol)
	pprofs.Get("/pprof/trace", pprof.Trace)
}

func (server *Server) Test(args GetArgs) {

}
