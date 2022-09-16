package pool

import (
	"context"
	"sync"
	"time"
)

// Pool 协程池
type Pool[T Goroutine] struct {
	tasks    []T           //任务缓存
	task     []chan T      //任务队列
	size     int           //任务处理数量 默认为0
	number   int           //协程数量
	cut      int           //切换
	catch    RunTimeErr    //协程错误处理器
	status   []bool        // 协程状态管理 检查当前协程是否是启动状态，用于处理协程内发生 panic 之后重新启动 Star 避免多开 协程
	mux      sync.Mutex    // 协程重启锁
	listener chan struct{} //监听器
	clear    context.CancelFunc
	rootCtx  context.Context // 协程上下文
	close    bool            // 协程池结束标识,close为false表示协程池处于运行中，反之挟持结束
}

// NewPool 初始化线程池
// number 协程数
// size 单协程任务数量
// ctx 上下文变量
func NewPool[T Goroutine](number, size int, ctx context.Context) *Pool[T] {
	p := new(Pool[T])
	p.number = number
	p.size = size
	p.task = make([]chan T, number)
	p.status = make([]bool, number)
	for i := 0; i < number; i++ {
		p.task[i] = make(chan T, p.size)
	}
	p.catch = poolPanic{}
	p.mux = sync.Mutex{}
	p.listener = make(chan struct{})

	//初始化根上下文
	cancel, cancelFunc := context.WithCancel(ctx)
	p.rootCtx = cancel
	p.clear = cancelFunc
	return p
}

// Start 启动协程池
func (pool *Pool[T]) Start() {
	pool.star()
	// 启动监听器 重启协程池
	go func() {
		for {
			select {
			// 重新启动协程池
			case <-pool.listener:
				go pool.star()

				//关闭重启监听
			case <-pool.rootCtx.Done():
				return
			}
		}
	}()
}

func (pool *Pool[T]) star() {
	// 启动 number 个协程
	for i := 0; i < pool.number; i++ {
		pool.mux.Lock()
		if !pool.status[i] && !pool.close {
			pool.status[i] = true
			ctx, _ := context.WithCancel(pool.rootCtx)
			go func(c context.Context, task chan T, id int, mx *sync.Mutex) {
				// 协程抛出错误之后 该协程会结束运行，为了保障 后来的任务能够正确的处理 需要重启当前的协程
				defer func(num int, rw *sync.Mutex) {
					// 标识当前协程 结束了,添加任务时候 校验
					rw.Lock()
					pool.status[num] = false
					rw.Unlock()
					// 重新运行
					pool.listener <- struct{}{}
				}(id, mx)
				defer pool.catch.Catch()
				for {
					select {

					// 读取 任务
					case t := <-task:
						t.Run(c)

					// 结束协程池
					case <-c.Done():
						mx.Lock()
						pool.close = true
						mx.Unlock()
						return
					}
				}
			}(ctx, pool.task[i], i, &pool.mux)
		}
		pool.mux.Unlock()
	}
}

// Add 添加任务
func (pool *Pool[T]) add(task T) {
	pool.mux.Lock()
	defer pool.mux.Unlock()
	// 更新轮询安排任务
	pool.next()
	// 查找一个正在运行的 协程池传递任务
	for {
		if pool.close {
			return
		}
		if !pool.check() {
			time.Sleep(time.Second * 2)
		}
		// 校验当前 协程池是否可用
		if pool.status[pool.cut] {
			pool.task[pool.cut] <- task
			return
		}
		// 不断的循环查找 直到找到为止
		pool.next()
	}
}

// 轮询器
func (pool *Pool[T]) next() {
	if pool.cut == pool.number-1 {
		pool.cut = 0
	} else {
		pool.cut++
	}
}

// 检查协程池是否存在可用协程
func (pool *Pool[T]) check() bool {
	for i := 0; i < pool.number; i++ {
		if pool.status[i] {
			return true
		}
	}
	return false
}

// Recover 自定义错误处理器
func (pool *Pool[T]) Recover(rec RunTimeErr) {
	pool.catch = rec
}

// Stop 关闭协程池
func (pool *Pool[T]) Stop() {
	pool.clear()
}
