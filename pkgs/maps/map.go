package maps

type Map[K comparable, V any] map[K]V

// New 创建一个map
func New[K comparable, V any]() Map[K, V] {
	m := make(Map[K, V])
	return m
}

// Put 存储一个k/v
func (receiver Map[K, V]) Put(key K, value V) {
	receiver[key] = value
}

// Get 获取一个k/v
func (receiver Map[K, V]) Get(key K) V {
	val := receiver[key]
	return val
}

// Delete 删除一个k/v
func (receiver Map[K, V]) Delete(key K) {
	delete(receiver, key)
}

// IsEmpty 判断map是否存储元素
func (receiver Map[K, V]) IsEmpty() bool {
	if len(receiver) == 0 {
		return true
	}
	return false
}
