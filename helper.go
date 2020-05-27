package assert

import (
	"fmt"
	"strconv"
)

func prepMsg(msg []interface{}, format string, args ...interface{}) []interface{} {
	return append(msg, "\n\n", fmt.Sprintf(format, args...))
}

func shell(i int) string {
	return "\x1B[" + strconv.Itoa(i) + "m"
}
