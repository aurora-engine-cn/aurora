# Aurora Web Framework

[![star](https://gitee.com/aurora-engine/aurora/badge/star.svg?theme=dark)](https://gitee.com/aurora-engine/aurora/stargazers)
[![Go Report Card](https://goreportcard.com/badge/gitee.com/aurora-engine/aurora)](https://goreportcard.com/report/gitee.com/aurora-engine/aurora)
[![License](https://img.shields.io/badge/license-apache--2.0-blue)](https://gitee.com/aurora-engine/aurora/blob/master/LICENSE)
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

## 项目案例参考
[community 仓库](https://gitee.com/aurora-engine/community)


## about the author

**作者:** Awen

**联系:** zhiwen_der@qq.com

## thanks

该框架参考了，HttpRouter 的字典树 方式来构建路由信息

![](https://camo.githubusercontent.com/5075c80d56620267702a3808e7a926ff51235b2ecd986441c092e3b6b821af83/68747470733a2f2f7265736f75726365732e6a6574627261696e732e636f6d2f73746f726167652f70726f64756374732f636f6d70616e792f6272616e642f6c6f676f732f6a625f6265616d2e737667)<br>
感谢 [JetBrains](https://www.jetbrains.com/) 支持了该开源项目

## copyright information

该项目签署了**Apache**授权许可，详情请参阅 [LICENSE](https://github.com/awensir/go-aurora/blob/main/LICENSE)
