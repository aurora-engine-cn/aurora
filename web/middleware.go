package web

// Middleware 中间件类型
type Middleware func(ctx Context) bool
