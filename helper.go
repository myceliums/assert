package assert

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"testing"
)

func prepMsg(msg []interface{}, format string, args ...interface{}) []interface{} {
	return append(msg, "\n\n", fmt.Sprintf(format, args...))
}

func shell(i int) string {
	return "\x1B[" + strconv.Itoa(i) + "m"
}

var ts = map[*runtime.Func]*testing.T{}

func t(a Assert) *testing.T {
	f := runtime.FuncForPC(reflect.ValueOf(a).Pointer())
	return ts[f]
}
