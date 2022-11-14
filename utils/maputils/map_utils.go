package maputils

// IsEmpty 判断map是否为空或是否存在元素
func IsEmpty[K comparable, V any](maps map[K]V) bool {
	return maps == nil || len(maps) == 0
}
