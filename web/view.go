package web

import (
	"embed"
	"html/template"
	"net/http"
	"path/filepath"
)

// ViewHandle 是整个服务器对视图渲染的核心函数,开发者实现改接口对需要展示的页面进行自定义处理
type ViewHandle func(string, string, embed.FS, http.ResponseWriter, Context)

func View(fullPath, relative string, static embed.FS, rew http.ResponseWriter, data Context) {
	var html *template.Template
	var err error
	if static != (embed.FS{}) {
		relative = filepath.ToSlash(relative)
		html, err = template.ParseFS(static, relative)
	} else {
		html, err = template.ParseFiles(fullPath)
	}
	if err != nil {
		panic(err)
	}
	err = html.Execute(rew, data)
	if err != nil {
		panic(err)
	}
}
