package example_test

import (
	"fmt"
	"gitee.com/aurora-engine/aurora"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"
)

type Aaa struct {
	Name string
	User *Bbb
}

type Bbb struct {
	Name string
	Mac  *Aaa
}

type Ccc struct {
	Name string
}

func TestAurora(t *testing.T) {
	a := aurora.NewAurora()
	/// 注册一个系统变量，类型为 *Ccc
	a.SysVariable(&Ccc{}, func(proxy *aurora.Proxy) interface{} {
		c := &Aaa{Name: "test"}
		return c
	})
	a.Get("/", func(ccc *Ccc, req *http.Request) {
		fmt.Println(ccc)
	})
	aurora.Run(a)
}

// TestGet Get请求测试
func TestGet(t *testing.T) {
	a := aurora.NewAurora()
	a.Get("/", func(r *http.Request, ctx aurora.Ctx) {

	})
	aurora.Run(a)
}

type Post struct {
	Name    string
	Age     int
	Gender  string
	Address []string
	Report  map[string]interface{}
}

// TestPost Post请求测试
func TestPost(t *testing.T) {
	a := aurora.NewAurora()
	a.Post("/post1", func(post Post) {
		fmt.Println(post)
	})

	a.Post("/post2", func(post *Post) {
		fmt.Println(post)
	})

	a.Post("/post3", func(post map[string]interface{}) {
		fmt.Println(post)
	})

	a.Post("/post4", func(post *map[string]interface{}) {
		fmt.Println(post)
	})

	aurora.Run(a)
}

// TestBenchmark 压测api
func TestBenchmark(t *testing.T) {
	a := aurora.NewAurora()
	a.Get("/", func() string {
		return ""
	})
	a.Get("/user/{id}", func(id interface{}) interface{} {
		return id
	})
	a.Post("/user", func() string {
		return ""
	})
	aurora.Run(a)
}

// 测试 接口
var api = make(map[string]interface{})
var chars = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u",
	"v", "w", "x", "y", "z"}

func init() {
	Api()
}
func Api() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10000; i++ {
		intn := rand.Intn(10)
		url := []string{}
		for j := 0; j < intn; j++ {
			l := rand.Intn(10)
			s := strings.Builder{}
			for k := l; k >= 0; k-- {
				s.WriteString(chars[rand.Intn(25)])
			}
			url = append(url, s.String())
		}
		join := strings.Join(url, "/")
		// 创建一个函数
		fun := func() {
			fmt.Println("/" + join)
		}
		// 注册到api中
		api["/"+join] = fun
	}
}

// TestRouter 路由测试
func TestRouter(t *testing.T) {
	a := aurora.NewAurora(aurora.Debug())
	for s, i := range api {
		a.Get(s, i)
	}
	go t.Run("http", TestHttp)
	aurora.Run(a)
}

func TestHttp(t *testing.T) {
	count := 0
	time.Sleep(time.Second * 10)
	for url, _ := range api {
		get, err := http.Get("http://localhost:8080" + url)
		if err != nil {
			log.Fatal(err)
			return
		}
		if get.Status == "200 OK" {
			fmt.Println(url, " 成功")
			count++
		} else {
			fmt.Println(url, " 失败")
		}
	}
	fmt.Println("成功访问了:", count, "个接口,一共 ", len(api), " 个接口")
}

// TestRESTFul 动态路由测试
func TestRESTFul(t *testing.T) {
	a := aurora.NewAurora()

	aurora.Run(a)
}
