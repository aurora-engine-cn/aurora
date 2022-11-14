package list

/*
	切片工具类
*/

type ArrayList[T any] []T

// Add 添加数据
func (al *ArrayList[T]) Add(e T) {
	*al = append(*al, e)
}

// Get 获取元素
func (al *ArrayList[T]) Get(i int) T {
	return (*al)[i]
}

// Remove 删除指定索引的数据
func (al *ArrayList[T]) Remove(index int) {
	if al == nil || index < 0 {
		return
	}
	*al = append((*al)[:index], (*al)[index+1:]...)
}

// Length 获取数据长度
func (al *ArrayList[T]) Length() int {
	if al == nil {
		return 0
	}
	return len(*al)
}
