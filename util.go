package zeroapi

import (
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

func funcName(i interface{}) string {
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	split1 := strings.Split(fn, ".")
	name := split1[len(split1)-1]
	split2 := strings.Split(name, "-")
	return split2[0]
}

func getValueByField(i interface{}, field string) (value interface{}) {
	v := reflect.ValueOf(i)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return
	}
	// 遍历结构体的所有字段
	for i := 0; i < v.NumField(); i++ {
		name := v.Type().Field(i).Name
		if name == "values" {
			// map[int32]interface{}
			for _, v := range v.Field(i).MapKeys() {
				if v.Int() == 2 {
					return v.Interface()
				}
			}
		}
	}
	f := v.FieldByName(field)
	if !f.IsValid() {
		return
	}
	value = f.Interface()
	return
}

func InterfaceToInt64(value interface{}) int64 {
	if value == nil {
		return 0
	}
	switch value.(type) {
	case int:
		return int64(value.(int))
	case int8:
		return int64(value.(int8))
	case int16:
		return int64(value.(int16))
	case int32:
		return int64(value.(int32))
	case int64:
		return value.(int64)
	case uint:
		return int64(value.(uint))
	case uint8:
		return int64(value.(uint8))
	case uint16:
		return int64(value.(uint16))
	case uint32:
		return int64(value.(uint32))
	case uint64:
		return int64(value.(uint64))
	case float32:
		return int64(value.(float32))
	case float64:
		return int64(value.(float64))
	case string:
		i, _ := strconv.ParseInt(value.(string), 10, 64)
		return i
	default:
		return 0
	}
}
