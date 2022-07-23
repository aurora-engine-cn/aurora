package aurora

import (
	"bytes"
	"context"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"reflect"
	"sync"
)

type Aurora struct {
	// 日志
	Log
	// 文件上传大小配置
	MaxMultipartMemory int64
	// name 服务名称
	name string
	// 服务器顶级上下文，通过此上下文可以跳过 go web 自带的子上下文去开启纯净的子go程，结束此上下文 web服务也将结束 <***>
	ctx context.Context
	// 结束上下文
	cancel func()
	// 主机信息
	host string
	// 服务端口号
	port string
	// 可指定的配置文件
	configpath string
	// 路由服务管理
	router *route
	// 项目根路径
	projectRoot string

	// 静态资源管理 默认为 root 目录
	resource string

	// 静态文件服务器接口
	fileService string

	// 常用的静态资源头
	resourceMapType ContentType

	// 自定义系统参数
	intrinsic map[string]Constructor

	// 分配代理实例
	proxyPool *sync.Pool
	// 分配路径构建
	pathPool *sync.Pool

	//Aurora 配置启动配置项
	opt []Option

	api map[string][]controlInfo
	// 各类配置项的存储，在初始化阶段预存了内置的配置项获取,可以通过api多这个配置项镜像添加或覆盖
	use map[interface{}]UseConfiguration
	// 最后初始化需要加载的配置项
	options []UseOption
	// ioc命名组件
	components []Component
	// 第三方组件管理容器
	component *ioc
	// 加载结构体作为处理器, 处理器并不会被注册到缓存容器中，处理器在启动期间会根据需要去缓存容器中寻找对应的依赖
	controllers []*reflect.Value
	// 配置实例，读取配置文件
	config Config
	// go web 原生服务器
	server *http.Server
	ln     net.Listener // web服务器监听,启动服务器时候初始化 <+++>  计划 使用 多路复用器
}

func NewAurora(option ...Option) *Aurora {
	//初始化日志
	logs := logrus.New()
	logs.SetFormatter(&Formatter{})
	logs.Out = os.Stdout
	a := &Aurora{
		port: "8080", //默认端口号
		router: &route{
			mx: &sync.Mutex{},
		},
		proxyPool: &sync.Pool{
			New: func() interface{} {
				return &Proxy{}
			},
		},
		pathPool: &sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
		server:    &http.Server{},
		resource:  "", //设定资源默认存储路径，需要连接项目更目录 和解析出来资源的路径，资源路径解析出来是没有前缀 “/” 的作为 resource属性，在其两边加上 斜杠
		use:       make(map[interface{}]UseConfiguration),
		component: newIoc(),
	}
	a.router.Aurora = a
	a.Log = logs
	projectRoot, _ := os.Getwd()
	a.projectRoot = projectRoot //初始化项目路径信息
	//执行配置项
	for _, opt := range option {
		opt(a)
	}
	middleware := new(Middleware)
	// 中间件配置项
	a.use[reflect.TypeOf(middleware).Elem()] = useMiddleware
	// 静态资源头配置项，主要设置可能不存在的资源头，或者过时的子资源
	a.use[reflect.TypeOf(ContentType{})] = useContentType
	// 命名组件
	a.use[reflect.TypeOf(Component{})] = useComponent
	// log 日志
	a.use[reflect.TypeOf(&logrus.Logger{})] = useLogrus
	// server
	a.use[reflect.TypeOf(&http.Server{})] = useServe
	a.viperConfig()
	return a
}

// Use 使用组件,把组件加载成为对应的配置
func (a *Aurora) Use(Configuration ...interface{}) {
	if Configuration == nil {
		return
	}
	if a.options == nil {
		a.options = make([]UseOption, 0)
	}
	var opt UseOption
	for _, u := range Configuration {
		rt := reflect.TypeOf(u)
		if useOption, b := a.use[rt]; b {
			opt = useOption(u)
			a.options = append(a.options, opt)
			continue
		}
		//检查是否是实现 Config配置接口
		if rt.Implements(reflect.TypeOf(new(Config)).Elem()) {
			opt = useConfig(u)
			a.options = append(a.options, opt)
			continue
		}

		//默认没有找到其他可配置项，把它当作处理器加载
		opt = useController(u)
		a.options = append(a.options, opt)
	}
}

// Run 启动服务器
func (a *Aurora) Run() error {
	err := a.componentInit()
	if err != nil {
		return err
	}
	var p, certFile, keyFile string
	if a.config != nil {
		p = a.config.GetString("aurora.server.port")
		certFile = a.config.GetString("aurora.server.tls.certFile")
		keyFile = a.config.GetString("aurora.server.tls.keyFile")
	}
	if p != "" {
		a.port = p
	}
	l, err := net.Listen("tcp", ":"+a.port)
	if err != nil {
		return err
	}
	a.ln = l
	if certFile != "" && keyFile != "" {
		return a.server.ServeTLS(l, certFile, keyFile)
	}
	return a.server.Serve(l) //启动服务器
}

func (a *Aurora) Root() string {
	return a.projectRoot
}
