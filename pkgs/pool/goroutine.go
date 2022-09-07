package pool

import "context"

// Goroutine 协程任务接口
type Goroutine interface {
	Run(ctx context.Context)
}

// RunTimeErr 协程池运行错误捕捉接口
type RunTimeErr interface {
	Catch()
}
