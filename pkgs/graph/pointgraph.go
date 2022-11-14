package graph

/*
	链表图结结构
*/

type Graph[T any] struct {
	points []*Point[T] // 点集合
	edges  []*Edge     // 边集合
}
