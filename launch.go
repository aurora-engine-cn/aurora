package aurora

// Application Web 应用程序接口
// 通过嵌套匿名 *Engine 实例来完成web服务的构建
type Application interface {
	// Use 加载配置
	Use(...interface{})

	// Server 服务器初始化
	// 依赖加载 等操作在这个函数中进行
	Server()

	// Router 路由加载函数
	Router()

	// ioc 容器启动 函数 该函数由 Aurora 实现
	start()
	// run 和 ioc 方法通过嵌套(继承 Aurora实例)
	run() error
}

// Run 启动服务器，启动阶段自动注册当前服务实例
func Run(app Application) error {
	// 注册当前服务
	app.Use(app)
	// 初始化 服务
	app.Server()
	// 启动ioc
	app.start()

	// 第三方库加载

	// 加载路由
	app.Router()
	// 运行服务器
	return app.run()
}
