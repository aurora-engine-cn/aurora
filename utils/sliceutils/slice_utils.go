package sliceutils

// IsEmpty 判断切片是否为空或是否存在元素
func IsEmpty[T any](slice []T) bool {
	return slice == nil || len(slice) == 0
}
