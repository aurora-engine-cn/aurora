# Aurora Web Framework

[![Stars](https://img.shields.io/github/stars/aurora-go/aurora)](https://gitee.com/aurora-engine/aurora)
[![Go Version](https://img.shields.io/badge/go1.16-version-orange)](https://golang.google.cn/)
[![Go Report Card](https://goreportcard.com/badge/gitee.com/aurora-engine/aurora)](https://goreportcard.com/report/gitee.com/aurora-engine/aurora)
[![License](https://img.shields.io/gitee.com/aurora-engine/aurora)](https://gitee.com/aurora-engine/aurora/blob/future/LICENSE)
<br>

Aurora 是用 Go(Golang) 编写的 Web 框架 ,将是 Golang 自诞生以来最好用的 Web 开发生产工具。路由处理灵活，集中式依赖管理，让项目代码结构更加优雅，专注于业务编码。
## go version
```
go1.16+
```

## start
```go
package main

import "gitee.com/aurora-engine/aurora"

func main() {
	//创建 实例
	a := aurora.NewAurora()
	//注册接口
	a.Get("/", func() {
		a.Info("hello web")
	})
	//启动服务器
	aurora.Run(a)
}
```

## document

[document](https://aurora-go.github.io)

## about the author

**作者:** Awen

**联系:** zhiwen_der@qq.com

## thanks

该框架参考了，HttpRouter 的字典树 方式来构建路由信息

感谢 [JetBrains](https://www.jetbrains.com/) 支持了该开源项目, 并提供了一年开发工具的支持

## copyright information

该项目签署了**Apache**授权许可，详情请参阅 [LICENSE](https://github.com/awensir/go-aurora/blob/main/LICENSE)
