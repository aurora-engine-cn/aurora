package web

// Component 命名组件
type Component map[string]any

// Constructor 用于加载 匿名组件的构造器类型
// Aurora 会执行这个函数 并把得到的变量加载到容器中
type Constructor func() any
