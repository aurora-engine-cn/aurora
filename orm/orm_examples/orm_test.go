package orm_examples

import (
	"database/sql"
	"fmt"
	"gitee.com/aurora-engine/aurora/orm"
	"gitee.com/aurora-engine/aurora/orm/sqlbuild"
	"gitee.com/aurora-engine/aurora/pkgs/list"
	"github.com/druidcaesa/ztool"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strings"
	"testing"
)

const dbUrl = "x:xxxx@tcp(x.x.x.x:3306)/x"

var open *sql.DB

func init() {
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		panic(err)
		return
	}
	open = db
}

func TestSQL(t *testing.T) {
	split := strings.Split("", " ")
	t.Log(split)
}

/*
column 定义规则 (同数据库定义顺序一致，之间用空格分开)
列名  列约束1 列约束2 列约束3
*/
type User struct {
	Id         string `column:"user_id"`
	Account    string `column:"user_account"`
	Name       string `column:"user_name"`
	Email      string `column:"user_email"`
	Password   string `column:"user_password"`
	Age        int    `column:"user_age"`
	Birthday   string `column:"user_birthday"`
	Head       string `column:"user_head_picture"`
	CreateTime string `column:"user_create_time"`
}

func (s *User) Table() string {
	return "comm_user"
}

func TestSql(t *testing.T) {
	s := sqlbuild.Sql()
	s.Select("s.name as nam ,age  ,time as t")
	s.Where("stu.name=1", "a.age=1 or ss.id='2' and nn=cc or bbb=bbb and ccc=ccc")
	//s.Or()
	s.Where("aa.bb='ssss'", "aa.bb=www")
	t.Log(s.String())
}

// 查询测试
func TestMapping_Selects(t *testing.T) {
	m := orm.CreateMapping[*User](open)
	stu := &User{
		Name: "awen",
	}
	l := list.ArrayList[*User]{}
	l = m.Selects(stu)
	for _, s := range l {
		fmt.Printf("%v\r\n", s)
	}
}

func TestMapping_SelectMap(t *testing.T) {
	m := orm.CreateMapping[*User](open)
	v := m.SelectMap(map[string]any{
		"user_name": "saber",
	})
	fmt.Printf("%v\n\r", v)
}

func TestMapping_SelectMaps(t *testing.T) {
	m := orm.CreateMapping[*User](open)
	l := list.ArrayList[*User]{}
	l = m.SelectMaps(map[string]any{
		"user_name": "saber",
	})
	for _, s := range l {
		fmt.Printf("%v\r\n", s)
	}
}

// 插入测试
func TestMapping_Insert(t *testing.T) {
	m := orm.CreateMapping[*User](open)
	uuid, err := ztool.IdUtils.SimpleUUID()
	if err != nil {
		t.Error(err.Error())
	}
	stu := &User{
		Id:         uuid,
		Account:    "12345678",
		Email:      "xxxxx@qq.com",
		Name:       "saber",
		Birthday:   ztool.DateUtils.Format(),
		CreateTime: ztool.DateUtils.Format(),
	}
	insert := m.Insert(stu)
	fmt.Println(insert)
}
func TestMapping_InsertMap(t *testing.T) {
	m := orm.CreateMapping[*User](open)
	uuid, err := ztool.IdUtils.SimpleUUID()
	if err != nil {
		t.Error(err.Error())
	}
	insert := m.InsertMap(map[string]any{
		"user_id":          uuid,
		"user_name":        "testMap",
		"user_account":     "123456789",
		"user_birthday":    ztool.DateUtils.Format(),
		"user_create_time": ztool.DateUtils.Format(),
	})
	fmt.Println(insert)
}

// 更新测试
func TestMapping_Update(t *testing.T) {
	m := orm.CreateMapping[*User](open)
	stu := &User{
		Account: "12345678",
		Email:   "xxxxx@qq.com",
		Name:    "saber",
	}
	value := &User{
		Account: "111111",
		Name:    "awen",
	}
	update := m.Update(stu, value)
	fmt.Println(update)
}

func TestMapping_UpdateMap(t *testing.T) {
	m := orm.CreateMapping[*User](open)
	c := map[string]any{
		"user_name":    "testMap",
		"user_account": "123456789",
	}
	v := map[string]any{
		"user_name": "testMapUpdate",
	}
	update := m.UpdateMap(c, v)
	fmt.Println(update)
}

// 删除测试
func TestMapping_Delete(t *testing.T) {
	m := orm.CreateMapping[*User](open)
	value := &User{
		Account: "111111",
		Name:    "awen",
	}
	d := m.Delete(value)
	fmt.Println(d)
}

func TestFunc(t *testing.T) {
	// 实现 函数逻辑
	swap := func(in []reflect.Value) []reflect.Value {
		return []reflect.Value{in[1], in[0]}
	}

	// 创建函数
	makeSwap := func(fptr any) {
		fn := reflect.ValueOf(fptr).Elem()
		// Make a function of the right type.
		v := reflect.MakeFunc(fn.Type(), swap)
		// Assign it to the value fn represents.
		fn.Set(v)
	}

	// Make and call a swap function for ints.
	var intSwap func(int, int) (int, int)
	makeSwap(&intSwap)
	fmt.Println(intSwap(0, 1))

	// Make and call a swap function for float64s.
	var floatSwap func(float64, float64) (float64, float64)
	makeSwap(&floatSwap)
	fmt.Println(floatSwap(2.72, 3.14))
}
