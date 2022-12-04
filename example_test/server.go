package example

import (
	"errors"
	"fmt"
	"gitee.com/aurora-engine/aurora"
	"gitee.com/aurora-engine/aurora/web"
	"net/http"
)

// Server 嵌套Aurora定义一个服务 实例
type Server struct {
	*aurora.Engine
}

type GetArgs struct {
	Name string
	Age  int `constraint:"check"`
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
	server.Constraint("check", func(value any) error {
		if value.(int) <= 0 {
			return errors.New("error value is 0")
		}
		return nil
	})
}

type T func(int)

func Test(a int) {

}

func Test2(t T) {

}

func (server *Server) Router() {
	// 添加 app 路由

	server.Post("/upload", func(data any, file *web.MultipartFile) {
		for _, headers := range file.File {
			err := file.SaveUploadedFile(headers[0], "W:\\code\\go\\framework\\aurora\\example_test\\"+headers[0].Filename)
			if err != nil {
				panic(err)
			}
		}
	})

	//server.Post("/user", func(name, age string) string {
	//	return ""
	//})
	//server.Get("/user/{id}", func(id string) string {
	//	return id
	//})
	//
	//pprofs := server.Group("/debug")
	//pprofs.Get("/pprof", pprof.Index)
	//pprofs.Get("/pprof/cmdline", pprof.Cmdline)
	//pprofs.Get("/pprof/profile", pprof.Profile)
	//pprofs.Get("/pprof/symbol", pprof.Symbol)
	//pprofs.Get("/pprof/trace", pprof.Trace)
}

func (server *Server) Test(args GetArgs) {

}
