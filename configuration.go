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

// Static web 静态资源配置
func Static(fs embed.FS) Option {
	return func(engine *Engine) {
		engine.router.Static(fs)
	}
}
