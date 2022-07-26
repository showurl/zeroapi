package zeroapi

import (
	"reflect"
	"runtime"
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
