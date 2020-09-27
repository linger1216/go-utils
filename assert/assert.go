package assert

import "fmt"

func Assert(condition bool, msg string, v ...interface{}) {
	if !condition {
		panic(fmt.Sprintf("assert: "+msg, v...))
	}
}

func Ensure(condition bool) {
	Assert(condition, "something wrong")
}
