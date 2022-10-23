package example

import (
	"gitee.com/aurora-engine/aurora"
)

// Server 嵌套Aurora定义一个服务 实例
type Server struct {
	*aurora.Engine
}

type GetArgs struct {
	Name string `empty:"false"`
	Age  int    `constraint:"isEmpty"`
}

func (server *Server) Server() {
	// 进行一下初始化操作，比如 控制器实例，全局中间件，全局变量，第三方依赖库等操作
}

func (server *Server) Router() {
	// 添加 app 路由
	//server.Get("/", func(name, age string) {
	//	fmt.Println("/")
	//})
	//server.Get("/abcde/aa/bb/{cc}", func(v string) {
	//	fmt.Printf("/abcde/aa/bb/%s\n", v)
	//})
	//
	//server.Get("/abcde/aa/bb/{cc}/{dd}", func(v, v1 string) {
	//	fmt.Printf("/abcde/aa/bb/%s/%s\n", v, v1)
	//})
	//
	//server.Get("/abc", func() {
	//	fmt.Printf("abc\n")
	//})
	//
	//server.Get("/abcs", func() {
	//	fmt.Printf("abcs\n")
	//})
	//
	//server.Get("/test", func(args GetArgs) {
	//
	//})
	//
	//server.Post("/user", func(name, age string) string {
	//	return ""
	//})
	//server.Get("/user/{id}", func(id string) string {
	//	return id
	//})
	//server.Get("/users/{id}", func(id string) string {
	//	return id
	//})
	//pprofs := server.Group("/debug")
	//pprofs.Get("/pprof", pprof.Index)
	//pprofs.Get("/pprof/cmdline", pprof.Cmdline)
	//pprofs.Get("/pprof/profile", pprof.Profile)
	//pprofs.Get("/pprof/symbol", pprof.Symbol)
	//pprofs.Get("/pprof/trace", pprof.Trace)

	server.Get("romane", func() {})
	server.Get("romanus", func() {})
	server.Get("romulus", func() {})
	server.Get("rubens", func() {})
	server.Get("ruber", func() {})
	server.Get("rubicon", func() {})
	server.Get("rubicundus", func() {})

}

func (server *Server) Test(args GetArgs) {

}
