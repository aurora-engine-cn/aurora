package aurora

import (
	"gitee.com/aurora-engine/aurora/core"
	"gitee.com/aurora-engine/aurora/web"
	"reflect"
)

// SysVariable 当前系统参数不支持泛型参数
func (engine *Engine) SysVariable(v interface{}, value web.System) {
	if v == nil || value == nil {
		return
	}
	rt := reflect.ValueOf(v)
	if engine.intrinsic == nil {
		engine.intrinsic = make(map[string]web.System)
	}
	key := core.BaseTypeKey(rt)
	engine.intrinsic[key] = value
}

// Aurora 系统变量
func req(ctx web.Context) interface{} {
	return ctx.Request()
}

func rew(ctx web.Context) interface{} {
	return ctx.Response()
}

func ctx(ctx web.Context) interface{} {
	return ctx
}

func file(ctx web.Context) interface{} {
	return ctx.MultipartFile()
}
