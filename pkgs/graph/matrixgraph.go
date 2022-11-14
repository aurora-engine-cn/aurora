package graph

import (
	"fmt"
	"gitee.com/aurora-engine/aurora/pkgs/queue"
	"gitee.com/aurora-engine/aurora/pkgs/stack"
)

/*
	邻接矩阵图结构
	矩阵图结构的 BFS 和 DFS 都和点存储的先后顺序相关，迭代的优先级均是按照顺序访问

*/

type MatrixGraph[T any] struct {
	points   []*Point[T]    //图点信息
	pointMap map[string]int //point id 查找 points索引
	mark     []bool         //访问标识符
	matrix   [][]*Edge      //构建存储图结构信息
}

// init 更具points的变化初始化 访问标记和矩阵图
func (receiver *MatrixGraph[T]) init() {
	l := len(receiver.points)
	receiver.mark = make([]bool, l)
	receiver.matrix = make([][]*Edge, l)
	for i := 0; i < l; i++ {
		receiver.matrix[i] = append(receiver.matrix[i], make([]*Edge, l)...)
	}
}

// Point 向途中添加一个点属性
func (receiver *MatrixGraph[T]) Point(id, name string, data T) {
	point := &Point[T]{Id: id, Name: name, DataInfo: data, JsonInfo: ""}
	receiver.Add(point)
}

func (receiver *MatrixGraph[T]) Add(point *Point[T]) {
	if receiver.points == nil {
		receiver.points = make([]*Point[T], 0)
	}
	l := len(receiver.points)
	if receiver.pointMap == nil {
		receiver.pointMap = make(map[string]int)
	}
	receiver.pointMap[point.Id] = l
	receiver.points = append(receiver.points, point)
	receiver.init()
}

// Drawing 画图 连线，连接两个点
// starId 起点id
// endId 结束点id
// wight 权重
func (receiver *MatrixGraph[T]) Drawing(starId, endId string, wight int) {
	s := receiver.pointMap[starId]
	e := receiver.pointMap[endId]
	if s == e {
		return
	}
	receiver.matrix[s][e] = &Edge{value: wight}
}

func (receiver *MatrixGraph[T]) DFS(id string) {
	defer receiver.refreshMark()
	i := receiver.pointMap[id]
	receiver.dfs(i)
}

// 递归 dfs 遍历
func (receiver *MatrixGraph[T]) dfs(index int) {
	for i := index; i < len(receiver.matrix); {
		if !receiver.mark[i] {
			// 当前节点访问
			receiver.mark[i] = true
			fmt.Println(receiver.points[i].Name)
		}
		// 访问相关的点
		for j := 0; j < len(receiver.matrix[i]); j++ {
			if receiver.matrix[i][j] != nil && !receiver.mark[j] {
				// 标记 第 j 个点被访问到了
				receiver.mark[j] = true
				fmt.Println(receiver.points[j].Name)
				// 继续从第 j 个点进行深入
				receiver.dfs(j)
			}
		}
		break
	}
}

func (receiver *MatrixGraph[T]) DFS_(id string) {
	defer receiver.refreshMark()
	i := receiver.pointMap[id]
	receiver.dfs_(i)
}

func (receiver *MatrixGraph[T]) dfs_(index int) {
	s := stack.New[int]()
	// 起始点入栈
	for i := index; i < len(receiver.matrix); i++ {
		if !receiver.mark[i] {
			// 当前节点访问
			receiver.mark[i] = true
			fmt.Println(receiver.points[i].Name)
		}
		for j := 0; j < len(receiver.matrix[i]); j++ {
			if receiver.matrix[i][j] != nil && !receiver.mark[j] {
				// 标记 第 j 个点被访问到了
				receiver.mark[j] = true
				fmt.Println(receiver.points[j].Name)
				// 继续从第 j 个点进行深入
				s.Push(j)
			}
		}
		if s.IsEmpty() {
			break
		}
		index = s.Popup()
	}

}

func (receiver *MatrixGraph[T]) BFS(id string) {
	defer receiver.refreshMark()
	i := receiver.pointMap[id]
	receiver.bfs(i)
}

func (receiver *MatrixGraph[T]) bfs(index int) {
	// 辅助队列
	q := queue.New[int]()
	for i := index; i < len(receiver.matrix); {
		for j := 0; j < len(receiver.matrix[i]); j++ {
			if !receiver.mark[i] {
				receiver.mark[i] = true
				fmt.Println(receiver.points[i].Name)
			}
			if receiver.matrix[i][j] != nil && !receiver.mark[j] {
				receiver.mark[j] = true
				fmt.Println(receiver.points[j].Name)
				q.EnQueue(j)
			}
		}
		if q.IsEmpty() {
			break
		}
		i = q.DeQueue()
	}
}

// 重置访问标识
func (receiver *MatrixGraph[T]) refreshMark() {
	receiver.mark = make([]bool, len(receiver.points))
}

func (receiver *MatrixGraph[T]) Print() {
	for i := range receiver.matrix {
		for j := range receiver.matrix[i] {
			edge := receiver.matrix[i][j]
			if edge != nil {
				fmt.Print(edge.value)
				continue
			}
			fmt.Print(0)

		}
		fmt.Println()
	}
}
