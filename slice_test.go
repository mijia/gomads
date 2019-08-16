package gomads

import (
	"fmt"
	"testing"
)

func TestBoxInts(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6}
	var b []string
	Box(a).Map(func(v int) string {
		return fmt.Sprintf("%d_str", v)
	}).Unbox(&b)

	fmt.Println(b)
}
