package gateway

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
