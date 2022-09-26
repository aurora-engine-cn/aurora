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
func (engine *Engine) ioc() {
	engine.Info("start component-dependent assembly")

	//加载uses配置项，配置项中可能存在加载ioc配置
	if engine.options != nil {
		for _, useOption := range engine.options {
			useOption(engine)
		}
	}

	// 加载 构造器 build 到 ioc 容器
	if engine.build != nil {
		for _, constructor := range engine.build {
			// 执行构造 生成组件放入到 ioc中
			c := constructor()
			engine.control(c)
		}
	}

	//加载 注册的依赖到 初级缓存
	if engine.components != nil {
		for _, component := range engine.components {
			for k, v := range component {
				if err := engine.component.putIn(k, v); err != nil {
					ErrorMsg(err)
				}
			}
		}
	}

	if engine.components != nil {
		for _, component := range engine.components {
			for k, v := range component {
				if err := engine.space.Put(k, v); err != nil {
					ErrorMsg(err)
				}
			}
		}
	}
	// 清空
	engine.components = nil
	//启动容器 ,给容器中的组件进行依赖初始化,容器加载出错 结束运行
	err := engine.space.Start()
	if err != nil {
		ErrorMsg(err)
	}

}

func (engine *Engine) StartRouter() {
	// 完成容器启动 ，这一步主要是针对于 属于controller处理器一部分进行操作，比如自动加载一些配置文件中的值
	// 该步骤仅对匿名的控制器组件产生效果，命名组件不处理
	engine.dependencyInjection()

	// 设置web服务的静态资源处理路径 默认初始化为 / 为 根路径
	// 此处的静态资源 不是作为文件服务器的支持 仅仅支持html资源的加载
	engine.resource = "/"

	engine.server.BaseContext = engine.baseContext //配置 上下文对象属性
	engine.router.defaultView = View               //初始化使用默认视图解析,aurora的视图解析是一个简单的实现，可以通过修改 a.Router.DefaultView 实现自定义的试图处理，框架最终调用此方法返回页面响应
	engine.server.Handler = engine.router          //设置默认路由器
	if engine.config != nil {                      //是否加载配置文件 覆盖配置项
		engine.Info("the configuration file is loaded successfully.")
		//加载配置文件中定义的 端口号
		port := engine.config.GetString("aurora.server.port")
		if port != "" {
			engine.port = port
		}
		//读取配置路径
		p := engine.config.GetString("aurora.resource")
		//构建路径拼接，此处在路径前后加上斜杠 用于静态资源的路径凭借方便
		if p != "" {
			if p[:1] != "/" {
				p = "/" + p
			}
			if p[len(p)-1:] != "/" {
				p = p + "/"
			}
			engine.resource = p
		}
		p = engine.config.GetString("aurora.server.file")
		engine.fileService = p
		engine.Info(fmt.Sprintf("server static resource root directory:%1s", engine.resource))
		name := engine.config.GetString("aurora.application.name")
		if name != "" {
			engine.name = name
			engine.Info("the service name is " + engine.name)
		}
	}

	//注册路由
	//if engine.api != nil {
	//	for method, infos := range engine.api {
	//		for _, info := range infos {
	//			engine.router.Register(method, info.path, info.control, info.middleware...)
	//		}
	//	}
	//	engine.api = nil
	//}
	engine.router.LoadCache()
}
