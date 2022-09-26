package aurora

import (
	"bytes"
	"context"
	"fmt"
	"gitee.com/aurora-engine/aurora/container"
	"gitee.com/aurora-engine/aurora/utils"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"sync"
)

var banner = " ,--.    __   _    _ .--.    .--.    _ .--.   ,--.\n`'_\\ :  [  | | |  [ `/'`\\] / .'`\\ \\ [ `/'`\\] `'_\\ :\n// | |,  | \\_/ |,  | |     | \\__. |  | |     // | |,\n\\'-;__/  '.__.'_/ [___]     '.__.'  [___]    \\'-;__/\n|          Aurora Web framework (v1.2.3)           |"

type Engine struct {
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
	router *Router
	// 项目根路径
	projectRoot string

	// 静态资源管理 默认为 root 目录
	resource string

	// 静态文件服务器接口
	fileService string

	// 自定义系统参数
	intrinsic map[string]Constructor

	// 分配代理实例
	proxyPool *sync.Pool
	// 分配路径构建
	pathPool *sync.Pool

	//Aurora 配置启动配置项
	opt []Option

	// 接口信息
	api map[string][]controlInfo

	// 各类配置项的存储，在初始化阶段预存了内置的配置项获取,可以通过api多这个配置项镜像添加或覆盖
	use map[interface{}]UseConfiguration

	// 最后初始化需要加载的配置项
	options []UseOption

	// ioc 命名组件
	components []Component

	// ioc 匿名组件
	build []Constructors

	// 第三方组件管理容器
	component *ioc
	space     *container.Space
	// 加载结构体作为处理器, 处理器并不会被注册到缓存容器中，处理器在启动期间会根据需要去缓存容器中寻找对应的依赖
	controllers []*reflect.Value

	// 配置实例，读取配置文件
	config Config

	// go web 原生服务器
	server *http.Server

	ln net.Listener

	// consul 治理中心
	consulCenter *ConsulCenter
}

func New(option ...Option) *Engine {
	engine := NewEngine()
	//初始化日志
	logs := logrus.New()
	logs.SetFormatter(&Formatter{})
	logs.Out = os.Stdout
	//初始化基本属性
	engine.Log = logs
	engine.printBanner()
	engine.Info(fmt.Sprintf("golang version :%1s", runtime.Version()))
	engine.control(engine)

	// 加载配置文件
	engine.viperConfig()
	engine.router = NewRoute(engine)
	// 加载 consul 配置
	engine.consul()

	// 执行配置项
	for _, opt := range option {
		opt(engine)
	}

	var middleware Middleware
	var constructors Constructors
	// 中间件配置项
	engine.use[reflect.TypeOf(middleware)] = useMiddleware
	// 匿名组件
	engine.use[reflect.TypeOf(constructors)] = useConstructors
	// 命名组件
	engine.use[reflect.TypeOf(Component{})] = useComponent
	// log 日志
	engine.use[reflect.TypeOf(&logrus.Logger{})] = useLogrus
	// server
	engine.use[reflect.TypeOf(&http.Server{})] = useServe
	return engine
}

// NewEngine 创建 Engine 基础配置
func NewEngine() *Engine {
	engine := &Engine{
		port: "8080", //默认端口号
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
		server:   &http.Server{},
		resource: "", //设定资源默认存储路径，需要连接项目更目录 和解析出来资源的路径，资源路径解析出来是没有前缀 “/” 的作为 resource属性，在其两边加上 斜杠
		use:      make(map[interface{}]UseConfiguration),
	}
	projectRoot, _ := os.Getwd()
	engine.projectRoot = projectRoot //初始化项目路径信息
	engine.component = newIoc(engine.Log)
	engine.space = container.NewSpace()
	// 初始化系统参数
	if engine.intrinsic == nil {
		engine.intrinsic = make(map[string]Constructor)
	}
	engine.intrinsic[utils.BaseTypeKey(reflect.ValueOf(new(http.Request)))] = req
	engine.intrinsic[utils.BaseTypeKey(reflect.ValueOf(new(http.ResponseWriter)).Elem())] = rew
	engine.intrinsic[utils.BaseTypeKey(reflect.ValueOf(new(Ctx)))] = ctx
	engine.intrinsic[utils.BaseTypeKey(reflect.ValueOf(new(MultipartFile)))] = file
	return engine
}

// Use 使用组件,把组件加载成为对应的配置
func (engine *Engine) Use(Configuration ...interface{}) {
	if Configuration == nil {
		return
	}
	if engine.options == nil {
		engine.options = make([]UseOption, 0)
	}
	var opt UseOption
	for _, u := range Configuration {
		rt := reflect.TypeOf(u)
		if useOption, b := engine.use[rt]; b {
			opt = useOption(u)
			engine.options = append(engine.options, opt)
			continue
		}
		//检查是否是实现 Config配置接口
		if rt.Implements(reflect.TypeOf(new(Config)).Elem()) {
			opt = useConfig(u)
			engine.options = append(engine.options, opt)
			continue
		}
		opt = useControl(u)
		engine.options = append(engine.options, opt)
	}
}

// Run 启动服务器
func (engine *Engine) run() error {
	// 启动路由
	engine.StartRouter()

	var p, certFile, keyFile string
	if engine.config != nil {
		p = engine.config.GetString("aurora.server.port")
		certFile = engine.config.GetString("aurora.server.tls.certFile")
		keyFile = engine.config.GetString("aurora.server.tls.keyFile")
	}
	if p != "" {
		engine.port = p
	}
	l, err := net.Listen("tcp", ":"+engine.port)
	if err != nil {
		return err
	}
	engine.ln = l
	if certFile != "" && keyFile != "" {

		return engine.server.ServeTLS(l, certFile, keyFile)
	}
	return engine.server.Serve(l)
}

// dependencyInjection Control 依赖加载
// controllers 属性中存储的都是 匿名组件类型
func (engine *Engine) dependencyInjection() {
	if engine.controllers == nil {
		return
	}
	engine.Info("Initialize load controller dependencies")
	l := len(engine.controllers)
	for i := 0; i < l; i++ {
		control := *engine.controllers[i]
		if control.Kind() == reflect.Ptr {
			control = control.Elem()
		}
		for j := 0; j < control.NumField(); j++ {
			field := control.Type().Field(j)
			//查询 value 属性 读取config中的基本属性
			if v, b := field.Tag.Lookup("value"); b {
				if v == "" {
					engine.Warn("value tag value is ''")
					continue
				}
				get := engine.config.Get(v)
				if get == nil {
					//如果查找结果大小等于0 则表示不存在
					continue
				}
				//把查询到的 value 初始化给指定字段
				err := utils.StarAssignment(control.Field(j), get)
				ErrorMsg(err)
			}
		}
	}
}

func (engine *Engine) printBanner() {
	fmt.Printf("%s\n\r", banner)
}

func (engine *Engine) Root() string {
	return engine.projectRoot
}
