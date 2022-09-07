package stack

type node[T comparable] struct {
	value T
	next  *node[T]
}
