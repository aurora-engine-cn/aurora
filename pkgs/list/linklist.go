package list

import "sync"

type LinkList[T comparable] struct {
	head *node[T]
	end  *node[T]
	iter *node[T]
	size int
	mx   sync.Mutex
}

// Add 顺序添加数据 T
func (receiver *LinkList[T]) Add(v T) {
	receiver.mx.Lock()
	defer receiver.mx.Unlock()
	n := &node[T]{value: v}
	if receiver.head == nil {
		receiver.iter = n
		receiver.head = n
		receiver.end = n
		receiver.size++
		return
	}
	receiver.end.next = n
	n.per = receiver.end
	receiver.end = n
	receiver.size++
}

func (receiver *LinkList[T]) Get(index int) T {
	var v T
	if receiver.size <= index {
		return v
	}
	n := 0
	p := receiver.iter
	for p != nil {
		if n < receiver.size && n == index {
			return p.value
		}
		p = p.next
		n++
	}
	return v
}

// Remove 顺序删除数据 T
func (receiver *LinkList[T]) Remove() {
	receiver.mx.Lock()
	defer receiver.mx.Unlock()
	if receiver.end != nil {
		n := receiver.end.per
		receiver.end = n
		receiver.size--
		return
	}
}
