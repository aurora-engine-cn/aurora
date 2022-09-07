package pages

// Page 分页数据
type Page[T any] struct {
	Count int64 `json:"count"` //总数
	Rows  []T   `json:"rows"`  //数据内容
}

func NewPage[T any](count int64, data ...T) *Page[T] {
	return &Page[T]{
		Count: count,
		Rows:  data,
	}
}

func (receiver *Page[T]) SetCount(count int64) {
	receiver.Count = count
}

func (receiver *Page[T]) SetData(data ...T) {
	receiver.Rows = data
}
