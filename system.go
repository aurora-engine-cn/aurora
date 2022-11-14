package aurora

import (
	"gitee.com/aurora-engine/aurora/core"
	"gitee.com/aurora-engine/aurora/web"
	"reflect"
)

// Variate 向 Engine 中注册一个内部的系统变量，
// value将返回一个和v相同的类型，否则panic
// 提供自定义系统变量注册，参数列表中的自定义类型需要严格匹配
func (engine *Engine) Variate(v any, value web.Variate) {
	if v == nil || value == nil {
		return
	}
	rt := reflect.ValueOf(v)
	if engine.intrinsic == nil {
		engine.intrinsic = make(map[string]web.Variate)
	}
	key := core.BaseTypeKey(rt)
	engine.intrinsic[key] = value
}

// Aurora 系统变量
func req(ctx web.Context) any {
	return ctx.Request()
}

func rew(ctx web.Context) any {
	return ctx.Response()
}

func ctx(ctx web.Context) any {
	return ctx
}

func file(ctx web.Context) any {
	return ctx.MultipartFile()
}
