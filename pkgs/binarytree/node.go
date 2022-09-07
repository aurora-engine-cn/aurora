package binarytree

type node[T any] struct {
	Data   T
	parent *node[T] //父节点
	left   *node[T] //左子节点
	right  *node[T] // 右子节点
}
