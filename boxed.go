package gomads

import (
	"fmt"
	"reflect"
)

// Boxed defines the boxed monads interface
type Boxed interface {
	Map(fmap interface{}) Boxed
	ConcatMap(fmap interface{}) Boxed
	Unbox(v interface{})
}

// Box will put the old container value into the monad box
func Box(v interface{}) Boxed {
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Slice:
		return newSliceBoxed(v)
	default:
		panic(fmt.Sprintf("no such boxed container support for %s", reflect.TypeOf(v)))
	}
}

func panicCondition(c bool, msg string) {
	if c {
		panic(msg)
	}
}
