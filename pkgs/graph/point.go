package graph

type Point[T any] struct {
	Id       string //点id
	Name     string //点名称
	DataInfo T      //点存储的信息
	JsonInfo string //json数据信息
}
