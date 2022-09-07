package orm

type ORM[T any] interface {
	// Insert 插入一个或者多个记录
	Insert(...T) int64
	InsertMap(...map[string]any) int64

	// Delete 删除满足条件的记录
	Delete(T) int
	DeleteMap(map[string]any) int

	// Update 更新指定内容
	Update(T, T) int64
	UpdateMap(map[string]any, map[string]any) int64

	Select(T) T
	SelectMap(map[string]any) T

	Selects(T) []T
	SelectMaps(map[string]any) []T
}
