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

var banner = " ,--.    __   _    _ .--.    .--.    _ .--.   ,--.\n`'_\\ :  [  | | |  [ `/'`\\] / .'`\\ \\ [ `/'`\\] `'_\\ :\n// | |,  | \\_/ |,  | |     | \\__. |  | |     // | |,\n\\'-;__/  '.__.'_/ [___]     '.__.'  [___]    \\'-;__/\n|          Aurora Web framework (v1.3.2)           |"

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
	intrinsic map[string]web.Variate

	//Aurora 配置启动配置项
	opt []Option

	// 各类配置项的存储，在初始化阶段预存了内置的配置项获取,可以通过api多这个配置项镜像添加或覆盖，use的配置项不支持泛型
	use map[string]useConfiguration

	// 最后初始化需要加载的配置项
	options []useOption

	// 命名组件
	components []web.Component

	// 匿名组件
	build []web.Constructor

	// 第三方组件管理容器
	space *container.Space

	// 配置实例，读取配置文件
	config web.Config

	// go app 原生服务器
	server *http.Server

	ln net.Listener
}

func New(option ...Option) *Engine {
	engine := NewEngine()
	engine.router = NewRoute(engine)
	// 执行配置项
	for _, opt := range option {
		opt(engine)
	}
	return engine
}

// NewEngine 创建 Engine 基础配置
// 初始化默认端口号
// 创建 http 服务
// 初始化默认日志
// 初始化系统参数列表
// 加载配置文件
func NewEngine() *Engine {
	engine := &Engine{
		port:     "8080", //默认端口号
		server:   &http.Server{},
		resource: "", //设定资源默认存储路径，需要连接项目更目录 和解析出来资源的路径，资源路径解析出来是没有前缀 “/” 的作为 resource属性，在其两边加上 斜杠
		use:      make(map[string]useConfiguration),
	}
	projectRoot, _ := os.Getwd()
	engine.projectRoot = projectRoot    //初始化项目路径信息
	engine.space = container.NewSpace() //初始化容器
	engine.space.Put("", engine)
	logs := logrus.New()
	logs.SetFormatter(&web.Formatter{})
	logs.Out = os.Stdout
	engine.Log = logs //初始化日志
	engine.printBanner()
	engine.Info(fmt.Sprintf("golang version :%1s", runtime.Version()))

	// 初始化 Use 配置
	key := ""
	var middleware web.Middleware
	var constructors web.Constructor
	// 中间件配置项
	key = core.TypeKey(middleware)
	engine.use[key] = useMiddleware
	// 匿名组件
	key = core.TypeKey(constructors)
	engine.use[key] = useConstructors
	// 命名组件
	key = core.TypeKey(web.Component{})
	engine.use[key] = useComponent
	// log 日志
	key = core.TypeKey(&logrus.Logger{})
	engine.use[key] = useLogrus
	// server
	key = core.TypeKey(&http.Server{})
	engine.use[key] = useServe

	// 初始化系统参数
	if engine.intrinsic == nil {
		engine.intrinsic = make(map[string]web.Variate)
	}
	key = core.BaseTypeKey(reflect.ValueOf(new(http.Request)))
	engine.intrinsic[key] = req
	key = core.BaseTypeKey(reflect.ValueOf(new(http.ResponseWriter)).Elem())
	engine.intrinsic[key] = rew
	key = core.BaseTypeKey(reflect.ValueOf(new(web.Context)))
	engine.intrinsic[key] = ctx
	key = core.BaseTypeKey(reflect.ValueOf(new(web.MultipartFile)))
	engine.intrinsic[key] = file
	// 加载配置文件
	engine.viperConfig()
	return engine
}

// NewRoute 创建并初始化 Router
func NewRoute(engine *Engine) *route.Router {
	router := route.New()
	router.MaxMultipartMemory = engine.MaxMultipartMemory
	router.Intrinsic = engine.intrinsic
	router.Log = engine.Log
	return router
}

// Use 使用组件,把组件加载成为对应的配置
func (engine *Engine) Use(Configuration ...any) {
	if Configuration == nil {
		return
	}
	if engine.options == nil {
		engine.options = make([]useOption, 0)
	}
	var opt useOption
	for _, u := range Configuration {
		key := core.TypeKey(u)
		if useOption, b := engine.use[key]; b {
			opt = useOption(u)
			engine.options = append(engine.options, opt)
			continue
		}
		opt = useControl(u)
		engine.options = append(engine.options, opt)
	}
}

// Verify 参数验证器
func (engine *Engine) Verify(tag string, fun web.Verify) {

}

// GetConfig 获取 Aurora 配置实例 对配置文件内容的读取都是协程安全的
func (engine *Engine) GetConfig() web.Config {
	return engine.config
}

func (engine *Engine) Root() string {
	return engine.projectRoot
}

// ViewHandle 修改默认视图解析接口
// Aurora 的路由树初始化默认使用的 Aurora 自己实现的视图解析
// 通过 该方法可以重新设置视图解析的逻辑处理，或者使用其他第三方的视图处理
// 现在的试图处理器处理方式比较局限，后续根据开发者需求进一步调整
func (engine *Engine) ViewHandle(v web.ViewHandle) {
	engine.router.DefaultView = v
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

// Run 启动服务器
func (engine *Engine) run() error {
	engine.server.BaseContext = engine.baseContext //配置 上下文对象属性
	engine.router.DefaultView = web.View           //初始化使用默认视图解析,aurora的视图解析是一个简单的实现，可以通过修改 a.Router.DefaultView 实现自定义的试图处理，框架最终调用此方法返回页面响应
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

// start 启动容器
func (engine *Engine) start() {
	engine.Info("start component-dependent assembly")

	//加载uses配置项
	if engine.options != nil {
		for _, useOption := range engine.options {
			useOption(engine)
		}
	}
	// 加载 构造器 build 到 容器
	if engine.build != nil {
		for _, constructor := range engine.build {
			// 执行构造 生成组件放入到 ioc中
			c := constructor()
			err := engine.space.Put("", c)
			if err != nil {
				panic(err)
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
	ErrorMsg(err, "Container initialization failed")
	engine.injection()
}

// injection  控制器依赖加载依赖加载,控制器的依赖加载实际在容器初始化阶段就已经完成
func (engine *Engine) injection() {

	// 获取容器中的主缓存
	Controllers := engine.space.Cache()
	for _, c := range Controllers {
		control := *c
		if control.Kind() == reflect.Ptr {
			control = control.Elem()
		}
		for j := 0; j < control.NumField(); j++ {
			field := control.Type().Field(j)
			//查询 value 属性 读取config中的基本属性
			if v, b := field.Tag.Lookup("value"); b && v != "" {
				get := engine.config.Get(v)
				if get == nil {
					//如果查找结果大小等于0 则表示不存在
					continue
				}
				//把查询到的 value 初始化给指定字段
				err := core.StarAssignment(control.Field(j), get)
				ErrorMsg(err)
			}
		}
	}
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
	if engine.config != nil { //是否加载配置文件 覆盖配置项
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

func (engine *Engine) printBanner() {
	fmt.Printf("%s\n\r", banner)
}

// baseContext 初始化 Aurora 顶级上下文
func (engine *Engine) baseContext(ln net.Listener) context.Context {
	c, f := context.WithCancel(context.TODO())
	//此处的保存在后续使用可能产生bug，情况未知
	engine.ctx = c
	engine.cancel = f
	engine.Info(fmt.Sprintf("the server successfully binds to the port:%s", engine.port))
	return c
}
