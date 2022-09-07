package maps

import (
	"encoding/json"
	"log"
)

// GetString 获取字符串值
func (receiver Map[K, V]) GetString(key K) string {
	if v, b := receiver[key]; b {
		marshal, err := json.Marshal(v)
		if err != nil {
			log.Println(err.Error())
			return ""
		}
		return string(marshal)
	} else {
		return ""
	}
}

func (receiver Map[K, V]) GetInt(key K) int {
	//if v, b := receiver[key]; b {
	//	//strconv.Atoi(v)
	//}
	return 0
}
