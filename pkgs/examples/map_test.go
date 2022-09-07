package examples

import (
	"fmt"
	"gitee.com/aurora-engine/pkgs/maps"
	"testing"
)

func TestMaps(t *testing.T) {
	m := maps.New[string, string]()
	m.Put("1", "a")
	for k, v := range m {
		fmt.Println("key:", k, " value:", v)
	}
	m.Delete("1")
	for k, v := range m {
		fmt.Println("key:", k, " value:", v)
	}
}
