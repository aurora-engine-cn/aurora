package example

import (
	"fmt"
	"gitee.com/aurora-engine/aurora"
	"testing"
)

func TestWebServer(t *testing.T) {

	err := aurora.Run(&Server{aurora.New(aurora.Debug())})
	if err != nil {
		fmt.Println(err)
		return
	}

}

func TestGin(t *testing.T) {

}

func BenchmarkRouter(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}
