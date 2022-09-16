package graph

type Point[T any] struct {
	Id        string         `json:"id"`        //点id
	Name      string         `json:"name"`      //点名称
	Type      string         `json:"type"`      //点类型
	Attribute map[string]any `json:"attribute"` //点属性
	DataInfo  T              `json:"dataInfo"`  //点存储的信息
	JsonInfo  string         `json:"jsonInfo"`  //json数据信息
}
