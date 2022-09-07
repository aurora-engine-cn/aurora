package pool

import "fmt"

type poolPanic struct {
}

func (p poolPanic) Catch() {
	if e := recover(); e != nil {
		// 协程异常处理
		fmt.Println(e)
	}
}
