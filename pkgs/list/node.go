package list

type node[T comparable] struct {
	value T
	per   *node[T] //前驱节点
	next  *node[T] //后继节点
}
