package aurora

import (
	"embed"
	"gitee.com/aurora-engine/aurora/web"
	"github.com/sirupsen/logrus"
	"reflect"
)

/*
	Aurora 配置项大全，New 启动阶段加载
	Use 阶段加载的配置会覆盖
*/

type Option func(*Engine)

// ConfigFilePath 指定 Aurora 加载配置文件位置
func ConfigFilePath(configPath string) Option {
	return func(a *Engine) {
		a.configpath = configPath
	}
}

// Config 指定 Aurora 的配置实例
func Config(config web.Config) Option {
	return func(engine *Engine) {
		engine.config = config
	}
}

// LoadConfig 加载配置文件数据
// 该方法只适用于 本地配置文件 embed 方式加载配置文件数据，初始化配置实例还是默认的
// 如果想要 第三方数据源 请使用 Config 方法替换掉 默认的配置实例
func LoadConfig(cnf []byte) Option {
	return func(engine *Engine) {
		engine.configFile = cnf
	}
}

// Debug 开启debug日志
func Debug() Option {
	return func(a *Engine) {
		of := reflect.ValueOf(a.Log)
		if of.Type().String() == reflect.ValueOf(&logrus.Logger{}).Type().String() {
			l := of.Interface()
			l.(*logrus.Logger).SetLevel(logrus.DebugLevel)
		}
	}
}

// MaxMultipartMemory 添加全局设置文件上传大小
func MaxMultipartMemory(size int64) Option {
	return func(engine *Engine) {
		engine.MaxMultipartMemory = size
	}
}

// Static web 静态资源配置
func Static(fs embed.FS) Option {
	return func(engine *Engine) {
		engine.router.Static(fs)
	}
}
