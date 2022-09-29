package web

// ControlInfo 用于存储在服务器启动之前注册的接口信息，需要在加载完配置项之后进行统一注册
type ControlInfo struct {
	Path       string
	Control    any
	Middleware []Middleware
}
