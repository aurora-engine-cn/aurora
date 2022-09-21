package queue

import "sync"

type Queue[T any] struct {
	head *node[T]
	end  *node[T]
	mx   sync.Mutex
	size int
}

func New[T comparable]() *Queue[T] {
	return new(Queue[T])
}

func (receiver *Queue[T]) EnQueue(value T) {
	receiver.mx.Lock()
	defer receiver.mx.Unlock()
	n := &node[T]{value: value}
	if receiver.head == nil {
		receiver.head = n
		receiver.end = n
		receiver.size++
		return
	}
	receiver.end.next = n
	receiver.end = n
	receiver.size++
}

func (receiver *Queue[T]) DeQueue() T {
	var v T
	receiver.mx.Lock()
	defer receiver.mx.Unlock()
	if receiver.head != nil {
		v = receiver.head.value
		receiver.head = receiver.head.next
		receiver.size--
		return v
	}
	return v
}

func (receiver *Queue[T]) IsEmpty() bool {
	if receiver.size == 0 {
		return true
	}
	return false
}

func (receiver *Queue[T]) Size() int {
	receiver.mx.Lock()
	defer receiver.mx.Unlock()
	return receiver.size
}
