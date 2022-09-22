package examples

import (
	"gitee.com/aurora-engine/aurora/pkgs/queue"
	"reflect"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	//q := queue.Queue[int]{}
	//q.EnQueue(1)
	//q.EnQueue(2)
	//fmt.Println(q.Size())
	//fmt.Println(q.DeQueue(), " ", q.Size())
	//fmt.Println(q.DeQueue(), " ", q.Size())
	//fmt.Println(q.DeQueue(), " ", q.Size())
}

func TestType(t *testing.T) {
	of := reflect.TypeOf(&queue.Queue[int]{})
	typeOf := reflect.TypeOf(queue.Queue[string]{})
	t.Log(of.Elem().PkgPath())
	t.Log(of.String())
	t.Log(of.PkgPath())
	t.Log(typeOf.String())

	t2 := reflect.TypeOf(time.Time{})
	t.Log(t2.PkgPath())
	t.Log(t2.String())
}
