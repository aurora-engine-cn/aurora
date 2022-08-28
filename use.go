package aurora

import (
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
		a.control(control)
	}
}

// useMiddleware Use的中间件注册
func useMiddleware(middleware interface{}) UseOption {
	return func(engine *Engine) {
		if m, b := middleware.(Middleware); !b {
			return
		} else {
			engine.router.use(m)
		}
	}
}

func useLogrus(log interface{}) UseOption {
	return func(engine *Engine) {
		engine.Log = log.(Log)
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

// useContentType 使用静态资源头
func useContentType(contentType interface{}) UseOption {
	return func(a *Engine) {
		if contentTypes, b := contentType.(ContentType); !b {
			return
		} else {
			for k, v := range contentTypes {
				a.resourceMapType[k] = v
			}
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
		if config, b := component.(Config); b {
			engine.config = config
		}
	}
}
