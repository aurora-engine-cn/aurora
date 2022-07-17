package aurora

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
)

// Component 组件加载类型
type Component map[string]interface{}

func (a *Aurora) componentInit() error {
	//加载uses配置项
	if a.options != nil {
		for _, useOption := range a.options {
			useOption(a)
		}
	}

	// 初始化系统参数
	if a.intrinsic == nil {
		a.intrinsic = make(map[string]Constructor)
	}
	a.intrinsic[reflect.TypeOf(&http.Request{}).String()] = req
	a.intrinsic[reflect.TypeOf(new(http.ResponseWriter)).Elem().String()] = rew
	a.intrinsic[reflect.TypeOf(Ctx{}).String()] = ctx
	a.intrinsic[reflect.TypeOf(&MultipartFile{}).String()] = file

	//初始化基本属性
	a.Info(fmt.Sprintf("golang version :%1s", runtime.Version()))
	a.loadResourceHead() //加载静态资源头

	//加载 注册的依赖到 初级缓存
	if a.components != nil {
		for _, component := range a.components {
			for k, v := range component {
				if err := a.component.putIn(k, v); err != nil {
					return err
				}
			}
		}
	}
	a.components = nil
	//启动容器 ,容器加载出错 结束运行
	err := a.component.start()
	if err != nil {
		return err
	}
	a.dependencyInjection()              // Ioc依赖注入初始化,第三方库的配置需要在此之前进行装配
	a.resource = "/"                     //默认初始化为 / 为 根路径
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
	return nil
}
