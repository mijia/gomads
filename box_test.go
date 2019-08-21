package gomads

import (
	"fmt"
	"strings"
	"testing"
)

func TestBoxInts(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6}
	var b []string
	Box(a).Map(func(v int) string {
		return fmt.Sprintf("%d_str", v)
	}).FlatMap(func(v string) []string {
		return []string{v, v, v}
	}).Unbox(&b)

	fmt.Println(b)
}

func TestBoxFunction(t *testing.T) {
	a := func(v int) bool {
		return v%2 == 0
	}
	var b func(v int) string
	Box(a).Map(func(v bool) string {
		if v {
			return "even"
		}
		return "odd"
	}).FlatMap(func(v string) func(x string) string {
		return func(x string) string {
			return "Tets"
		}
	}).Unbox(&b)
	fmt.Println(b(12))
	fmt.Println(b(13))
}

func TestBoxChannel(t *testing.T) {
	ch := make(chan int)
	var outChan chan string
	Box(ch).Map(func(v int) string {
		return fmt.Sprintf("%d_str_test", v)
	}).FlatMap(func(v string) chan string {
		out := make(chan string)
		go func() {
			sp := strings.Split(v, "_")
			for _, p := range sp {
				out <- p
			}
			close(out)
		}()
		return out
	}).Unbox(&outChan)

	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}()

	for v := range outChan {
		fmt.Println(v)
	}
}
