package examples

import (
	"fmt"
	"gitee.com/aurora-engine/pkgs/queue"
	"testing"
)

func TestQueue(t *testing.T) {
	q := queue.Queue[int]{}
	q.EnQueue(1)
	q.EnQueue(2)
	fmt.Println(q.Size())
	fmt.Println(q.DeQueue(), " ", q.Size())
	fmt.Println(q.DeQueue(), " ", q.Size())
	fmt.Println(q.DeQueue(), " ", q.Size())
}
