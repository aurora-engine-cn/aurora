package jstruct

import (
	jsoniter "github.com/json-iterator/go"
	"strconv"
)

type JSONStruct map[string]interface{}

func (js JSONStruct) Json() string {
	marshal, err := jsoniter.Marshal(js)
	if err != nil {
		panic(err)
	}
	return string(marshal)
}

func (js JSONStruct) Get(key string) interface{} {
	if v, b := js[key]; b {
		return v
	}
	return nil
}
func (js JSONStruct) GetString(key string) (string, bool) {
	if v, b := js[key]; b {
		var value string
		switch v.(type) {
		case string:
			value = v.(string)
		case int:
			i := v.(int)
			value = strconv.Itoa(i)
		case float64:
			f := v.(float64)
			value = strconv.FormatFloat(f, 'f', 'f', 64)
		case bool:
			b2 := v.(bool)
			value = strconv.FormatBool(b2)
		}
		return value, b
	}
	return "", false
}
func (js JSONStruct) GetInt(key string) (int, bool) {
	if v, b := js[key]; b {
		var value int
		switch v.(type) {
		case string:
			s := v.(string)
			atoi, err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			}
			value = atoi
		case int:
			value = v.(int)

		case float64:
			value = (int)(v.(float64))

		case bool:
			b2 := v.(bool)
			if b2 {
				value = 1

			}
			value = 0
		}
		return value, b
	}
	return 0, false
}
func (js JSONStruct) GetFloat64(key string) (float64, bool) {
	if v, b := js[key]; b {
		var value float64
		switch v.(type) {
		case string:
			s := v.(string)
			float, err := strconv.ParseFloat(s, 64)
			if err != nil {
				panic(err)
			}
			value = float
		case int:
			value = float64(v.(int))

		case float64:
			value = v.(float64)

		case bool:
			b2 := v.(bool)
			if b2 {
				value = 1
			}
			value = 0
		}
		return value, b
	}
	return 0, false
}
func (js JSONStruct) GetBool(key string) (bool, bool) {
	if v, b := js[key]; b {
		var value bool
		switch v.(type) {
		case string:
			s := v.(string)
			bol, err := strconv.ParseBool(s)
			if err != nil {
				panic(err)
			}
			value = bol
		case int:
			i := float64(v.(int))
			value = i != 0
		case float64:
			f := v.(float64)
			value = f != 0
		case bool:
			value = v.(bool)
		}
		return value, b
	}
	return false, false
}
func (js JSONStruct) GetMap(key string) (map[string]interface{}, bool) {
	if v, b := js[key]; b {
		var value map[string]interface{}
		switch v.(type) {
		case map[string]interface{}:
			value = v.(map[string]interface{})
		default:
			b = false
		}
		return value, b
	}
	return nil, false
}

func (js JSONStruct) GetSlice(key string) []interface{} {
	if v, b := js[key]; b {
		if value, f := v.([]interface{}); f {
			return value
		}
	}
	return nil
}
func (js JSONStruct) GetStringSlice(key string) []string {
	if v, b := js[key]; b {
		switch v.(type) {
		case []string:
			s := v.([]string)
			return s
		case []int:
			arr := v.([]int)
			s := make([]string, len(arr))
			for i, value := range arr {
				itoa := strconv.Itoa(value)
				s[i] = itoa
			}
			return s
		case []float64:
			arr := v.([]float64)
			s := make([]string, len(arr))
			for i, value := range arr {
				float := strconv.FormatFloat(value, 'f', 'f', 64)
				s[i] = float
			}
			return s
		case []bool:
			arr := v.([]bool)
			s := make([]string, len(arr))
			for i, value := range arr {
				formatBool := strconv.FormatBool(value)
				s[i] = formatBool
			}
			return s
		}
	}
	return nil
}
func (js JSONStruct) GetIntSlice(key string) []int {
	if v, b := js[key]; b {
		var value []int
		switch v.(type) {
		case []string:
			arr := v.([]string)
			value = make([]int, len(arr))
			for i, vs := range arr {
				atoi, err := strconv.Atoi(vs)
				if err != nil {
					panic(err)
				}
				value[i] = atoi
			}
		case []int:
			value = v.([]int)
		case []float64:
			arr := v.([]float64)
			value = make([]int, len(arr))
			for i, vf := range arr {
				value[i] = int(vf)
			}
		case []bool:
			arr := v.([]bool)
			value = make([]int, len(arr))
			for i, vb := range arr {
				if vb {
					value[i] = 1
				} else {
					value[i] = 0
				}
			}
		default:

		}
		return value
	}
	return nil
}
func (js JSONStruct) GetFloat64Slice(key string) []float64 {
	if v, b := js[key]; b {
		var value []float64
		switch v.(type) {
		case []string:
			arr := v.([]string)
			value = make([]float64, len(arr))
			for i, vs := range arr {
				atoi, err := strconv.ParseFloat(vs, 64)
				if err != nil {
					panic(err)
				}
				value[i] = atoi
			}
		case []int:
			arr := v.([]int)
			value = make([]float64, len(arr))
			for i, va := range arr {
				value[i] = float64(va)
			}
		case []float64:
			value = v.([]float64)
		case []bool:
			arr := v.([]bool)
			value = make([]float64, len(arr))
			for i, vb := range arr {
				if vb {
					value[i] = 1
				} else {
					value[i] = 0
				}
			}
		default:

		}
		return value
	}
	return nil
}
