# gomads

Brain training with Golang monadic practice

Please checkout the article: [https://awalterschulze.github.io/blog/post/monads-for-goprogrammers/](https://awalterschulze.github.io/blog/post/monads-for-goprogrammers/)

## Error monadic handling

```go
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
fmt.Println(b([]byte(`"1"`)))     // odd <nil>
fmt.Println(b([]byte(`"2"`)))     // even <nil>
fmt.Println(b([]byte(`"0"`)))     // "" 0 is not even nor odd
fmt.Println(b([]byte(`"test"`)))  // "" strconv.Atoi: parsing "test": invalid syntax
```
