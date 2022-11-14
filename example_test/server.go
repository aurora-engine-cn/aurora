package example

import (
	"fmt"
	"gitee.com/aurora-engine/aurora"
	"gitee.com/aurora-engine/aurora/web"
	"net/http"
	"net/http/pprof"
)

// Server 嵌套Aurora定义一个服务 实例
type Server struct {
	*aurora.Engine
}

type GetArgs struct {
	Name string `empty:"false" value:""`
	Age  int    `constraint:"check"`
}

func Recover() web.Recover {
	return func(w http.ResponseWriter) {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}
}

func (server *Server) Server() {
	server.Use(Recover())
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
