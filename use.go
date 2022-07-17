package aurora

import (
	"net/http"
)

// UseOption 配置项 对 *Aurora 的指定属性进行 配置
type UseOption func(*Aurora)

type UseConfiguration func(interface{}) UseOption

// useController Use的 处理器注册
func useController(control interface{}) UseOption {
	return func(a *Aurora) {
		a.control(control)
	}
}

// useMiddleware Use的中间件注册
func useMiddleware(middleware interface{}) UseOption {
	return func(a *Aurora) {
		if m, b := middleware.(Middleware); !b {
			return
		} else {
			a.router.use(m)
		}
	}
}

func useLogrus(log interface{}) UseOption {
	return func(a *Aurora) {
		a.Log = log.(Log)
	}
}

// useServe 使用自定义的 serve 实例
func useServe(server interface{}) UseOption {
	return func(a *Aurora) {
		if server == nil {
			return
		}
		if s, b := server.(*http.Server); !b {
			return
		} else {
			a.server = s
		}

	}
}

// useContentType 使用静态资源头
func useContentType(contentType interface{}) UseOption {
	return func(a *Aurora) {
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
	return func(a *Aurora) {
		if c, b := component.(Component); b {
			if a.components == nil {
				a.components = make([]Component, 0)
				a.components = append(a.components, c)
				return
			}
			a.components = append(a.components, c)
		}
	}
}

// useConfig 使用自定义viper配置
func useConfig(component interface{}) UseOption {
	return func(a *Aurora) {
		if component == nil {
			return
		}
		if config, b := component.(Config); b {
			a.config = config
		}
	}
}
