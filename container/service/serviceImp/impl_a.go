package serviceImp

type Aaa struct {
	Name string
}

func (a *Aaa) Get(name string) string {
	return a.Name
}
