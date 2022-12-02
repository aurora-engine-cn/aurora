package container

import (
	"gitee.com/aurora-engine/aurora/container/service/serviceImp"
	"testing"
)

type Aaa struct {
	Bbb
	*Ccc
	Name string
}

type Bbb struct {
	Name string
	*Ccc
}

type Ccc struct {
	Name string
	*Ddd
}

type Ddd struct {
	Name string
	*Fff
}

type Eee struct {
	Name string
}

type Fff struct {
	Name string
	*Aaa
	*Eee
}

func TestSpace_Start(t *testing.T) {
	aaa := &Aaa{Name: "aaa"}
	bbb := &Bbb{Name: "bbb"}
	ccc := &Ccc{Name: "ccc"}
	ddd := &Ddd{Name: "ddd"}
	eee := &Eee{Name: "eee"}
	//fff := &Fff{Name: "fff"}
	space := NewSpace()
	space.Put("", aaa)
	space.Put("", bbb)
	space.Put("", ccc)
	space.Put("", ddd)
	space.Put("", eee)
	//space.Put("", fff)
	err := space.Start()
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log()
}

func TestImpl(t *testing.T) {
	aaa := &serviceImp.Aaa{Name: "aaa"}
	bbb := &serviceImp.Bbb{}
	space := NewSpace()
	space.Put("A", aaa)
	space.Put("", bbb)
	err := space.Start()
	if err != nil {
		t.Error(err.Error())
		return
	}
	aaa.Name = "bbb"
	get := bbb.Get()
	t.Log(get)
}
