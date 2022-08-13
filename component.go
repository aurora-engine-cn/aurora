package aurora

import (
	"fmt"
)

// Component 命名组件
type Component map[string]interface{}

// Constructors 用于加载 匿名组件的构造器类型
// Aurora 会执行这个函数 并把得到的变量加载到 ioc 容器中
type Constructors func() interface{}

// StartIoc 启动容器
func (a *Aurora) ioc() {
	a.Info("start component-dependent assembly")

	//加载uses配置项，配置项中可能存在加载ioc配置
	if a.options != nil {
		for _, useOption := range a.options {
			useOption(a)
		}
	}

	// 加载 构造器 build 到 ioc 容器
	if a.build != nil {
		for _, constructor := range a.build {
			// 执行构造 生成组件放入到 ioc中
			c := constructor()

			a.control(c)
		}
	}

	//加载 注册的依赖到 初级缓存
	if a.components != nil {
		for _, component := range a.components {
			for k, v := range component {
				if err := a.component.putIn(k, v); err != nil {
					ErrorMsg(err)
				}
			}
		}
	}
	// 清空
	a.components = nil
	//启动容器 ,给容器中的组件进行依赖初始化,容器加载出错 结束运行
	err := a.component.start()
	if err != nil {
		ErrorMsg(err)
	}

}

func (a *Aurora) startRouter() {
	// 完成容器启动 ，这一步主要是针对于 属于controller处理器一部分进行操作，比如自动加载一些配置文件中的值
	// 该步骤仅对匿名的控制器组件产生效果，命名组件不处理
	a.dependencyInjection()

	// 设置web服务的静态资源处理路径 默认初始化为 / 为 根路径
	// 此处的静态资源 不是作为文件服务器的支持 仅仅支持html资源的加载
	a.resource = "/"

	a.server.BaseContext = a.baseContext //配置 上下文对象属性
	a.router.defaultView = a             //初始化使用默认视图解析,aurora的视图解析是一个简单的实现，可以通过修改 a.Router.DefaultView 实现自定义的试图处理，框架最终调用此方法返回页面响应
	a.server.Handler = a                 //设置默认路由器
	if a.config != nil {                 //是否加载配置文件 覆盖配置项
		a.Info("the configuration file is loaded successfully.")
		//加载配置文件中定义的 端口号
		port := a.config.GetString("aurora.server.port")
		if port != "" {
			a.port = port
		}
		//读取配置路径
		p := a.config.GetString("aurora.resource")
		//构建路径拼接，此处在路径前后加上斜杠 用于静态资源的路径凭借方便
		if p != "" {
			if p[:1] != "/" {
				p = "/" + p
			}
			if p[len(p)-1:] != "/" {
				p = p + "/"
			}
			a.resource = p
		}
		p = a.config.GetString("aurora.server.file")
		a.fileService = p
		a.Log.Info(fmt.Sprintf("server static resource root directory:%1s", a.resource))
		name := a.config.GetString("aurora.application.name")
		if name != "" {
			a.name = name
			a.Info("the service name is " + a.name)
		}
	}
	//注册路由
	if a.api != nil {
		for method, infos := range a.api {
			for _, info := range infos {
				a.router.addRoute(method, info.path, info.control, info.middleware...)
			}
		}
		a.api = nil
	}
}
