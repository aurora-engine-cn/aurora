package aurora

import (
	"context"
	"fmt"
	"html/template"
	"net"
	"net/http"
)

// ViewHandle 修改默认视图解析接口
// Aurora 的路由树初始化默认使用的 Aurora 自己实现的视图解析
// 通过 该方法可以重新设置视图解析的逻辑处理，或者使用其他第三方的视图处理
// 现在的试图处理器处理方式比较局限，后续根据开发者需求进一步调整
func (a *Aurora) viewHandle(v views) {
	a.router.defaultView = v
}

// View 默认视图解析
// html: 需要被处理的静态资源绝对路径信息
// data: 是一个可传递的数据
func (a *Aurora) view(html string, rew http.ResponseWriter, data interface{}) {

	parseFiles, err := template.ParseFiles(html)
	if err != nil {
		a.Error(err.Error())
		return
	}
	err = parseFiles.Execute(rew, data)
	if err != nil {
		a.Error(err.Error())
		return
	}
}

// baseContext 初始化 Aurora 顶级上下文
func (a *Aurora) baseContext(ln net.Listener) context.Context {
	c, f := context.WithCancel(context.TODO())
	//此处的保存在后续使用可能产生bug，情况未知
	a.ctx = c
	a.cancel = f
	a.Info(fmt.Sprintf("the server successfully binds to the port:%s", a.port))
	return c
}
