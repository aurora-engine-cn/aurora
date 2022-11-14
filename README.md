# Aurora Web Framework
![logo](https://img-1252940994.cos.ap-nanjing.myqcloud.com/logo.png)<br>
[![star](https://gitee.com/aurora-engine/aurora/badge/star.svg?theme=dark)](https://gitee.com/aurora-engine/aurora/stargazers)
[![Go Report Card](https://goreportcard.com/badge/gitee.com/aurora-engine/aurora)](https://goreportcard.com/report/gitee.com/aurora-engine/aurora)
[![License](https://img.shields.io/badge/license-apache--2.0-blue)](https://gitee.com/aurora-engine/aurora/blob/master/LICENSE)
<br>

Aurora 是用 Go(Golang) 编写的 Web 框架 ,将是 Golang 自诞生以来最好用的 Web 开发生产工具。为了打造更友好的Go Web开发环境，框架的项目结构和开发习惯借鉴了著名框架 `Gin` 和 `Spring Boot` 的开发设计，框架设计采用了 `Gin` 等 Go 框架的 HTTP 注册方式和责任链调用的中间件处理，
同时结合了 `Spring Boot` 框架的请求参数解析和响应方式 。 简单且强大的同时保障了代码结构的优雅。将是 Golang 自诞生以来最好用的 Web 开发生产工具，
项目托管平台已经转移到 Gitee， 交流群:836414068， 如果您觉得 aurora 不错，或者对您有帮助，请赏颗星吧！
## Go 版本
```
go1.19
```

## 快速开始

## 导入
```go
import (
    "gitee.com/aurora-engine/aurora"
)

```

创建一个结构体，嵌套一个匿名`*aurora.Engine` 实例 完成对服务器的创建
```go
// Server 嵌套Aurora定义一个服务 实例
type Server struct {
    *aurora.Engine
}
```
实现 `aurora.Application` 接口中的两个方法,`Server()` 和 `Router()`
```go
func (server *Server) Server() {
	// 进行一下初始化操作，比如 控制器实例，全局中间件，全局变量，第三方依赖库等操作
}

func (server *Server) Router() {
	// 添加 app 路由
	server.Get("/", func() string {
		return "hello world"
	})
}
```

通过执行器启动web服务即可
```go
err := aurora.Run(&Server{aurora.New(aurora.Debug())})
if err != nil {
	fmt.Println(err)
	return
}
```

## 文档
有关更多的使用操作请查看 [最新在线文档](https://go-aurora-engine.github.io)

## 关于作者

**作者:** Awen

**联系:** zhiwen_der@qq.com

## 致谢
![](https://camo.githubusercontent.com/5075c80d56620267702a3808e7a926ff51235b2ecd986441c092e3b6b821af83/68747470733a2f2f7265736f75726365732e6a6574627261696e732e636f6d2f73746f726167652f70726f64756374732f636f6d70616e792f6272616e642f6c6f676f732f6a625f6265616d2e737667)<br>
感谢 [JetBrains](https://www.jetbrains.com/) 支持了该开源项目

## 版权信息

该项目签署了**Apache**授权许可，详情请参阅 [LICENSE](https://gitee.com/aurora-engine/aurora/blob/new_dev/LICENSE)
