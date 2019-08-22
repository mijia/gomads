package gomads

import (
	"encoding/json"
	"fmt"
	"strconv"
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
	a := func(v int) (bool, error) {
		if v == 0 {
			return false, fmt.Errorf("0 is not even nor odd")
		}
		return v%2 == 0, nil
	}
	var b func(v int) (string, error)
	Box(a).Map(func(v bool) string {
		if v {
			return "even"
		}
		return "odd"
	}).FlatMap(func(v string) (string, error) {
		if v == "even" {
			return "", fmt.Errorf("we don't process the even target")
		}
		return v, nil
	}).Unbox(&b)
	fmt.Println(b(0))
	fmt.Println(b(12))
	fmt.Println(b(13))
}

func TestBoxFunction2(t *testing.T) {
	a := func(data []byte) (s string, err error) {
		err = json.Unmarshal(data, &s)
		return
	}
	var b func(data []byte) (int, error)
	Box(a).FlatMap(strconv.Atoi).Unbox(&b)
	fmt.Println(b([]byte(`"1"`)))
	fmt.Println(b([]byte(`"test"`)))
}

func TestComposeErrors(t *testing.T) {
	var b func(data []byte) (string, error)
	ComposeErrors(
		func(data []byte) (s string, err error) {
			err = json.Unmarshal(data, &s)
			return
		},
		strconv.Atoi,
		func(i int) (bool, error) {
			if i == 0 {
				return false, fmt.Errorf("0 is not even nor odd")
			}
			return i%2 == 0, nil
		},
		func(isEven bool) string {
			if isEven {
				return "even"
			}
			return "odd"
		},
	).Unbox(&b)
	fmt.Println(b([]byte(`"1"`)))
	fmt.Println(b([]byte(`"2"`)))
	fmt.Println(b([]byte(`"0"`)))
	fmt.Println(b([]byte(`"test"`)))
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
