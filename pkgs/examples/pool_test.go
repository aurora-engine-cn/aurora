package examples

import (
	"context"
	"fmt"
	"gitee.com/aurora-engine/aurora/pkgs/pool"
	"testing"
	"time"
)

type Task struct {
	Id int
}

func (receiver Task) Run(ctx context.Context) {
	if receiver.Id%3 == 0 || receiver.Id%5 == 0 {
		panic(fmt.Sprintf("panic Task %d \n", receiver.Id))
	}
	fmt.Printf("Task %d \n", receiver.Id)
}

func TestPool(t *testing.T) {
	newPool := pool.NewPool[Task](10, 10, context.TODO())
	newPool.Start()
	for i := 4; i < 100; i++ {
		//time.Sleep(50 * time.Millisecond)
		go func(id int) {
			newPool.Execute(Task{id})
		}(i)
	}
	time.Sleep(10000 * time.Millisecond)
}
