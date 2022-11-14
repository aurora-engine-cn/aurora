package list

type List[T comparable] interface {
	Add(T)
	Get(int) *T
	Remove(int)
	Length() int
}
