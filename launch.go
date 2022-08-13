package aurora

type Application interface {
	Use(...interface{})
	Run() error
}

func Run(app Application) {
	// 注册自身
	app.Use(app)
	err := app.Run()
	ErrorMsg(err)
}

// Use 提供一个全局的注册器，把参数 components 加载到 Aurora实例中
func Use(app Application, components ...Component) {
	if components == nil {
		return
	}
	for _, component := range components {
		app.Use(component)
	}
}
