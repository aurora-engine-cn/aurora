package pool

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Pool 协程池
type Pool[T Goroutine] struct {
	task     []chan T           //每个通道表示一个协程
	size     int                //每个协程 候选任务处理数量 默认为0
	number   int                //协程数量
	cut      int                //切换，用于添加任务选择 协程的轮询器变量
	iterMux  sync.Mutex         // 保证轮询器的迭代在 多个协程中并发安全
	catch    RunTimeErr         //协程错误处理器
	status   []bool             // 协程状态管理 检查当前协程是否是启动状态，用于处理协程内发生 panic 之后重新启动 Star 避免多开 协程
	mux      sync.Mutex         // 协程重启锁
	listener chan struct{}      //监听器
	clear    context.CancelFunc //结束子协程
	rootCtx  context.Context    // 协程上下文
	close    bool               // 协程池结束标识,close为false表示协程池处于运行中，反之挟持结束
}

// NewPool 初始化线程池
// number 协程数
// size 单协程任务数量
// ctx 上下文变量
func NewPool[T Goroutine](number, size int, ctx context.Context) *Pool[T] {
	p := new(Pool[T])
	p.number = number
	p.size = size
	if size == 0 {
		p.size = 1
	}
	p.task = make([]chan T, number)
	p.status = make([]bool, number)
	for i := 0; i < number; i++ {
		p.task[i] = make(chan T, p.size)
	}
	p.catch = poolPanic{}
	p.mux = sync.Mutex{}
	p.iterMux = sync.Mutex{}
	p.listener = make(chan struct{})
	p.close = true

	var c context.Context
	var cancelFunc context.CancelFunc
	//初始化根上下文
	if ctx != nil {
		c, cancelFunc = context.WithCancel(ctx)
	} else {
		c, cancelFunc = context.WithCancel(context.Background())
	}
	p.rootCtx = c
	p.clear = cancelFunc
	return p
}

// Start 启动协程池
func (pool *Pool[T]) Start() {
	pool.close = false
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
				// 关闭所有通道
				for i := 0; i < len(pool.task); i++ {
					close(pool.task[i])
				}
				return
			}
		}
	}()
}

// Execute 执行任务
func (pool *Pool[T]) Execute(task T) {
	pool.add(task)
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
				defer pool.reload(id)
				defer pool.catch.Catch()
				for {
					select {

					// 读取 任务
					case t, ok := <-task:
						//协程池关闭后 将不继续处理任务
						if ok {
							fmt.Printf("goroutine %d running .. ", id)
							t.Run(c)
						}
					// 结束协程池
					case <-c.Done():
						pool.close = true
						return
					}
				}
			}(ctx, pool.task[i], i, &pool.mux)
		}
		pool.mux.Unlock()
	}
}

// add 添加任务
func (pool *Pool[T]) add(task T) {
	if pool.close {
		return
	}
	// 更新轮询安排任务
	pool.next()
	// 查找一个正在运行的 协程池传递任务
	for {
		if !pool.check() {
			// 全都不可用时 等待2秒的回复时间
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
	pool.iterMux.Lock()
	pool.iterMux.Unlock()
	if pool.cut == pool.number-1 {
		pool.cut = 0
	} else {
		pool.cut++
	}
}

// 检查协程池是否存在可用协程,只要存在可用状态的协程 就返回true
func (pool *Pool[T]) check() bool {
	for i := 0; i < pool.number; i++ {
		if pool.status[i] {
			return true
		}
	}
	return false
}

// 协程任务运行中出现panic级别错误用于回复 协程运行
func (pool *Pool[T]) reload(num int) {
	pool.mux.Lock()
	defer pool.mux.Unlock()
	if pool.close {
		return
	}
	// 标识当前协程 结束了,添加任务时候 校验
	pool.status[num] = false
	// 重新运行
	pool.listener <- struct{}{}
}

// Recover 自定义错误处理器
func (pool *Pool[T]) Recover(rec RunTimeErr) {
	pool.catch = rec
}

// Stop 关闭协程池
// 执行关闭后 出现panic的协程将不会重启，正在运行的 协程执行玩当前任务也将结束
func (pool *Pool[T]) Stop() {
	pool.mux.Lock()
	defer pool.mux.Unlock()
	pool.clear()
}
