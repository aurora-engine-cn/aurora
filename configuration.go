package aurora

import (
	"gitee.com/aurora-engine/aurora/web"
	"github.com/sirupsen/logrus"
	"reflect"
)

/*
	Aurora 配置项大全，New 启动阶段加载
	Use 阶段加载的配置会覆盖
*/

type Option func(*Engine)

// ConfigFile 指定Aurora加载配置文件
func ConfigFile(configPath string) Option {
	return func(a *Engine) {
		a.configpath = configPath
	}
}

func Config(config web.Config) Option {
	return func(engine *Engine) {
		engine.config = config
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
