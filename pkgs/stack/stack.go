package stack

import "sync"

type Stack[T any] struct {
	top  *node[T]
	mx   sync.Mutex
	size int
}

func New[T any]() *Stack[T] {
	return new(Stack[T])
}

func (receiver *Stack[T]) Push(value T) {
	receiver.mx.Lock()
	defer receiver.mx.Unlock()
	n := &node[T]{value: value}
	n.next = receiver.top
	receiver.top = n
	receiver.size++
}

func (receiver *Stack[T]) Popup() T {
	receiver.mx.Lock()
	defer receiver.mx.Unlock()
	var v T
	if receiver.top != nil {
		v = receiver.top.value
		receiver.top = receiver.top.next
		receiver.size--
		return v
	}
	return v
}

// Size 获取栈大小
func (receiver *Stack[T]) Size() int {
	receiver.mx.Lock()
	defer receiver.mx.Unlock()
	return receiver.size
}

func (receiver *Stack[T]) IsEmpty() bool {
	receiver.mx.Lock()
	defer receiver.mx.Unlock()
	if receiver.top == nil {
		return true
	}
	return false
}
