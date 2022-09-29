package aurora

import (
	"context"
	"fmt"
	"gitee.com/aurora-engine/aurora/container"
	"gitee.com/aurora-engine/aurora/route"
	"gitee.com/aurora-engine/aurora/utils"
	"gitee.com/aurora-engine/aurora/web"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
)

var banner = " ,--.    __   _    _ .--.    .--.    _ .--.   ,--.\n`'_\\ :  [  | | |  [ `/'`\\] / .'`\\ \\ [ `/'`\\] `'_\\ :\n// | |,  | \\_/ |,  | |     | \\__. |  | |     // | |,\n\\'-;__/  '.__.'_/ [___]     '.__.'  [___]    \\'-;__/\n|          Aurora Web framework (v1.2.3)           |"

type Engine struct {
	// 日志
	web.Log
	// 文件上传大小配置
	MaxMultipartMemory int64
	// name 服务名称
	name string
	// 服务器顶级上下文，通过此上下文可以跳过 go app 自带的子上下文去开启纯净的子go程，结束此上下文 web服务也将结束 <***>
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
	router *route.Router
	// 项目根路径
	projectRoot string

	// 静态资源管理 默认为 root 目录
	resource string

	// 静态文件服务器接口
	fileService string

	// 自定义系统参数
	intrinsic map[string]route.Constructor

	//Aurora 配置启动配置项
	opt []Option

	// 各类配置项的存储，在初始化阶段预存了内置的配置项获取,可以通过api多这个配置项镜像添加或覆盖
	use map[interface{}]UseConfiguration

	// 最后初始化需要加载的配置项
	options []UseOption

	// 命名组件
	components []Component

	// 匿名组件
	build []Constructors

	// 第三方组件管理容器
	space *container.Space

	// 加载结构体作为处理器, 处理器并不会被注册到缓存容器中，处理器在启动期间会根据需要去缓存容器中寻找对应的依赖
	controllers []*reflect.Value

	// 配置实例，读取配置文件
	config web.Config

	// go app 原生服务器
	server *http.Server

	ln net.Listener

	// consul 治理中心
	consulCenter *ConsulCenter
}

func New(option ...Option) *Engine {
	engine := NewEngine()
	engine.router = NewRoute(engine)
	// 加载 consul 配置
	engine.consul()

	// 执行配置项
	for _, opt := range option {
		opt(engine)
	}

	var middleware web.Middleware
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
		port:     "8080", //默认端口号
		server:   &http.Server{},
		resource: "", //设定资源默认存储路径，需要连接项目更目录 和解析出来资源的路径，资源路径解析出来是没有前缀 “/” 的作为 resource属性，在其两边加上 斜杠
		use:      make(map[interface{}]UseConfiguration),
	}
	projectRoot, _ := os.Getwd()
	engine.projectRoot = projectRoot    //初始化项目路径信息
	engine.space = container.NewSpace() //初始化容器
	logs := logrus.New()
	logs.SetFormatter(&Formatter{})
	logs.Out = os.Stdout
	engine.Log = logs //初始化日志
	engine.printBanner()
	engine.Info(fmt.Sprintf("golang version :%1s", runtime.Version()))
	engine.space.Put("", engine) // 把自己注册到容器中
	// 初始化系统参数
	if engine.intrinsic == nil {
		engine.intrinsic = make(map[string]route.Constructor)
	}
	engine.intrinsic[utils.BaseTypeKey(reflect.ValueOf(new(http.Request)))] = req
	engine.intrinsic[utils.BaseTypeKey(reflect.ValueOf(new(http.ResponseWriter)).Elem())] = rew
	engine.intrinsic[utils.BaseTypeKey(reflect.ValueOf(new(web.Context)))] = ctx
	engine.intrinsic[utils.BaseTypeKey(reflect.ValueOf(new(web.MultipartFile)))] = file
	// 加载配置文件
	engine.viperConfig()
	return engine
}

func NewRoute(engine *Engine) *route.Router {
	router := route.New()
	router.MaxMultipartMemory = engine.MaxMultipartMemory
	router.Intrinsic = engine.intrinsic
	router.Log = engine.Log
	return router
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
		if rt.Implements(reflect.TypeOf(new(web.Config)).Elem()) {
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
	engine.server.BaseContext = engine.baseContext //配置 上下文对象属性
	engine.router.DefaultView = engine.View        //初始化使用默认视图解析,aurora的视图解析是一个简单的实现，可以通过修改 a.Router.DefaultView 实现自定义的试图处理，框架最终调用此方法返回页面响应
	engine.server.Handler = engine.router          //设置默认路由器
	engine.router.LoadCache()                      //加载接口
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

func (engine *Engine) printBanner() {
	fmt.Printf("%s\n\r", banner)
}

func (engine *Engine) Root() string {
	return engine.projectRoot
}

func ErrorMsg(err error, msg ...string) {
	if err != nil {
		if msg == nil {
			msg = []string{"Error"}
		}
		emsg := fmt.Errorf("%s : %s", strings.Join(msg, ""), err.Error())
		panic(emsg)
	}
}
