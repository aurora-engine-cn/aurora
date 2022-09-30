package web

import (
	"html/template"
	"net/http"
)

// ViewHandle 是整个服务器对视图渲染的核心函数,开发者实现改接口对需要展示的页面进行自定义处理
type ViewHandle func(string, http.ResponseWriter, Context)

func View(html string, rew http.ResponseWriter, data Context) {
	parseFiles, err := template.ParseFiles(html)
	if err != nil {
		panic(err)
	}
	err = parseFiles.Execute(rew, data)
	if err != nil {
		panic(err)
	}
}
