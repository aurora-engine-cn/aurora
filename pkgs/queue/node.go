package queue

/*
	队列节点
*/

type node[T comparable] struct {
	value T
	next  *node[T]
}
