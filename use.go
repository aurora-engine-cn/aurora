package aurora

import (
	"gitee.com/aurora-engine/aurora/web"
	"net/http"
)

// UseOption 配置项 对 *Aurora 的指定属性进行 配置
type UseOption func(engine *Engine)

type UseConfiguration func(interface{}) UseOption

// useController Use的 处理器注册
func useConstructors(control interface{}) UseOption {
	return func(engine *Engine) {
		if constructors, b := control.(Constructors); b {
			if engine.build == nil {
				engine.build = make([]Constructors, 0)
			}
			engine.build = append(engine.build, constructors)
		}
	}
}

// useControl
func useControl(control interface{}) UseOption {
	return func(a *Engine) {
		err := a.space.Put("", control)
		if err != nil {
			panic(err)
		}
	}
}

// useMiddleware Use 的中间件注册
func useMiddleware(middleware interface{}) UseOption {
	return func(engine *Engine) {
		if m, b := middleware.(web.Middleware); !b {
			return
		} else {
			engine.router.Use(m)
		}
	}
}

func useLogrus(log interface{}) UseOption {
	return func(engine *Engine) {
		engine.Log = log.(web.Log)
	}
}

// useServe 使用自定义的 serve 实例
func useServe(server interface{}) UseOption {
	return func(engine *Engine) {
		if server == nil {
			return
		}
		if s, b := server.(*http.Server); !b {
			return
		} else {
			engine.server = s
		}

	}
}

// useComponent 添加到容器
func useComponent(component interface{}) UseOption {
	return func(engine *Engine) {
		if c, b := component.(Component); b {
			if engine.components == nil {
				engine.components = make([]Component, 0)
				engine.components = append(engine.components, c)
				return
			}
			engine.components = append(engine.components, c)
		}
	}
}

// useConfig 使用自定义viper配置
func useConfig(component interface{}) UseOption {
	return func(engine *Engine) {
		if component == nil {
			return
		}
		if config, b := component.(web.Config); b {
			engine.config = config
		}
	}
}
