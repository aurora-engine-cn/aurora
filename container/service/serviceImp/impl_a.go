package serviceImp

type Aaa struct {
	Name string
}

func (a *Aaa) Get() string {
	return a.Name
}
