package aurora

import (
	"context"
	"fmt"
	"gitee.com/aurora-engine/aurora/container"
	"gitee.com/aurora-engine/aurora/core"
	"gitee.com/aurora-engine/aurora/route"
	"gitee.com/aurora-engine/aurora/web"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
)
const (
	yml  = "application.yml"
	yaml = "application.yaml"
)
var banner = " ,--.    __   _    _ .--.    .--.    _ .--.   ,--.\n`'_\\ :  [  | | |  [ `/'`\\] / .'`\\ \\ [ `/'`\\] `'_\\ :\n// | |,  | \\_/ |,  | |     | \\__. |  | |     // | |,\n\\'-;__/  '.__.'_/ [___]     '.__.'  [___]    \\'-;__/\n|          Aurora Web framework (v1.3.1)           |"

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
	intrinsic map[string]web.System

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
	logs.SetFormatter(&web.Formatter{})
	logs.Out = os.Stdout
	engine.Log = logs //初始化日志
	engine.printBanner()
	engine.Info(fmt.Sprintf("golang version :%1s", runtime.Version()))
	// 初始化系统参数
	if engine.intrinsic == nil {
		engine.intrinsic = make(map[string]web.System)
	}
	engine.intrinsic[core.BaseTypeKey(reflect.ValueOf(new(http.Request)))] = req
	engine.intrinsic[core.BaseTypeKey(reflect.ValueOf(new(http.ResponseWriter)).Elem())] = rew
	engine.intrinsic[core.BaseTypeKey(reflect.ValueOf(new(web.Context)))] = ctx
	engine.intrinsic[core.BaseTypeKey(reflect.ValueOf(new(web.MultipartFile)))] = file
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

// viperConfig 配置并加载 application.yml 配置文件
func (engine *Engine) viperConfig() {
	var ConfigPath string
	var err error
	if engine.configpath == "" {
		// 扫描配置文件
		filepath.WalkDir(engine.projectRoot, func(p string, d fs.DirEntry, err error) error {
			//找到配置及文件,基于根目录优先加载最外层的application.yml
			if !d.IsDir() && (strings.HasSuffix(p, yml) || (strings.HasSuffix(p, yaml))) && ConfigPath == "" {
				//修复 项目加载配置覆盖，检索项目配置文件，避免内层同名配置文件覆盖外层，这个情况可能发生在 开发者把两个go mod 项目嵌套在一起，导致配置被覆盖
				//此处校验，根据检索的更路径，只加载最外层的配置文件
				ConfigPath = p
			}
			return nil
		})
	} else {
		ConfigPath = engine.configpath
	}
	if ConfigPath == "" {
		engine.config = &web.ConfigCenter{Viper: viper.New(), RWMutex: &sync.RWMutex{}}
		return
	}
	if engine.config == nil {
		// 用户没有提供 配置项 则创建默认的配置处理
		cnf := &web.ConfigCenter{
			viper.New(),
			&sync.RWMutex{},
		}
		cnf.SetConfigFile(ConfigPath)
		err = cnf.ReadInConfig()
		ErrorMsg(err)
		engine.config = cnf
	}
	// 加载基础配置
	if engine.config != nil {                      //是否加载配置文件 覆盖配置项
		engine.Info("the configuration file is loaded successfully.")
		// 读取web服务端口号配置
		port := engine.config.GetString("aurora.server.port")
		if port != "" {
			engine.port = port
		}
		// 读取静态资源配置路径
		engine.resource = "/"
		p := engine.config.GetString("aurora.resource")
		// 构建路径拼接，此处在路径前后加上斜杠 用于静态资源的路径凭借方便
		if p != "" {
			if p[:1] != "/" {
				p = "/" + p
			}
			if p[len(p)-1:] != "/" {
				p = p + "/"
			}
			engine.resource = p
		}
		// 读取文件服务配置
		p = engine.config.GetString("aurora.server.file")
		engine.fileService = p
		engine.Info(fmt.Sprintf("server static resource root directory:%1s", engine.resource))
		// 读取服务名称
		name := engine.config.GetString("aurora.application.name")
		if name != "" {
			engine.name = name
			engine.Info("the service name is " + engine.name)
		}
	}
}

// GetConfig 获取 Aurora 配置实例 对配置文件内容的读取都是协程安全的
func (engine *Engine) GetConfig() web.Config {
	return engine.config
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
