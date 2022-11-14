package aurora

import (
	"gitee.com/aurora-engine/aurora/web"
	"net/http"
)

// UseOption 配置项 对 *Aurora 的指定属性进行 配置
type useOption func(engine *Engine)

type useConfiguration func(interface{}) useOption

// useController Use的 处理器注册
func useConstructors(control any) useOption {
	return func(engine *Engine) {
		if constructors, b := control.(web.Constructor); b {
			if engine.build == nil {
				engine.build = make([]web.Constructor, 0)
			}
			engine.build = append(engine.build, constructors)
		}
	}
}

// useControl
func useControl(control any) useOption {
	return func(a *Engine) {
		err := a.space.Put("", control)
		if err != nil {
			panic(err)
		}
	}
}

// useMiddleware Use 的中间件注册
func useMiddleware(middleware any) useOption {
	return func(engine *Engine) {
		if m, b := middleware.(web.Middleware); !b {
			return
		} else {
			engine.router.Use(m)
		}
	}
}

func useLogrus(log any) useOption {
	return func(engine *Engine) {
		engine.Log = log.(web.Log)
	}
}

// useServe 使用自定义的 serve 实例
func useServe(server any) useOption {
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

func useRecover(Recover any) useOption {
	return func(engine *Engine) {
		if recovers, b := Recover.(web.Recover); b {
			engine.router.Recover(recovers)
		}
	}
}

// useComponent 添加到容器
func useComponent(component any) useOption {
	return func(engine *Engine) {
		if c, b := component.(web.Component); b {
			if engine.components == nil {
				engine.components = make([]web.Component, 0)
				engine.components = append(engine.components, c)
				return
			}
			engine.components = append(engine.components, c)
		}
	}
}

// useConfig 使用自定义viper配置
func useConfig(component any) useOption {
	return func(engine *Engine) {
		if component == nil {
			return
		}
		if config, b := component.(web.Config); b {
			engine.config = config
		}
	}
}
