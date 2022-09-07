package list

type List[T comparable] interface {
	Add(...T)
	Get(int) *T
	Delete(int)
	Length() int
}
