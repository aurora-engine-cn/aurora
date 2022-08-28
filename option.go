package aurora

import (
	"github.com/sirupsen/logrus"
	"reflect"
)

/*
	Aurora 配置项大全
*/

type Option func(*Engine)

// ConfigFile 指定Aurora加载配置文件
func ConfigFile(configPath string) Option {
	return func(a *Engine) {
		a.configpath = configPath
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
