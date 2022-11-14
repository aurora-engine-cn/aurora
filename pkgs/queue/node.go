package queue

/*
	队列节点
*/

type node[T any] struct {
	value T
	next  *node[T]
}
