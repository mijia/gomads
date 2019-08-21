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
