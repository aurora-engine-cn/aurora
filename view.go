package aurora

import (
	"context"
	"fmt"
	"gitee.com/aurora-engine/aurora/route"
	"gitee.com/aurora-engine/aurora/web"
	"html/template"
	"net"
	"net/http"
)

// ViewHandle 修改默认视图解析接口
// Aurora 的路由树初始化默认使用的 Aurora 自己实现的视图解析
// 通过 该方法可以重新设置视图解析的逻辑处理，或者使用其他第三方的视图处理
// 现在的试图处理器处理方式比较局限，后续根据开发者需求进一步调整
func (engine *Engine) ViewHandle(v route.ViewHandle) {
	engine.router.DefaultView = v
}

// View 默认视图解析
// html: 需要被处理的静态资源绝对路径信息
// data: 是一个可传递的数据
func (engine *Engine) View(html string, rew http.ResponseWriter, data web.Context) {
	parseFiles, err := template.ParseFiles(html)
	ErrorMsg(err)
	err = parseFiles.Execute(rew, data)
	ErrorMsg(err)
}

// baseContext 初始化 Aurora 顶级上下文
func (engine *Engine) baseContext(ln net.Listener) context.Context {
	c, f := context.WithCancel(context.TODO())
	//此处的保存在后续使用可能产生bug，情况未知
	engine.ctx = c
	engine.cancel = f
	engine.Info(fmt.Sprintf("the server successfully binds to the port:%s", engine.port))
	return c
}
