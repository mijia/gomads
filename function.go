package gomads

import "reflect"

type _FunctionBoxed struct {
	T reflect.Type
	V reflect.Value
}

func newFunctionBoxed(v interface{}) Boxed {
	fn, err := baseType(reflect.TypeOf(v), reflect.Func)
	panicCondition(err != nil, "gomads: functionBox got no function ("+reflect.TypeOf(v).String()+")")
	b := _FunctionBoxed{
		T: fn,
		V: reflect.ValueOf(v),
	}
	return &b
}

func (b *_FunctionBoxed) Map(fmap interface{}) Boxed {
	ft := reflect.TypeOf(fmap)
	panicCondition(ft.Kind() != reflect.Func, "gomads: Map (not function)")
	panicCondition(b.T.NumOut() != ft.NumIn(), "gomads: Map (function has no type tranformation)")
	for i := 0; i < b.T.NumOut(); i++ {
		panicCondition(b.T.Out(i) != ft.In(i), "gomads: Map (function has no type transformation)")
	}
	inTypes := make([]reflect.Type, 0, b.T.NumIn())
	for i := 0; i < b.T.NumIn(); i++ {
		inTypes = append(inTypes, b.T.In(i))
	}
	outTypes := make([]reflect.Type, 0, ft.NumOut())
	for i := 0; i < ft.NumOut(); i++ {
		outTypes = append(outTypes, ft.Out(i))
	}
	newFuncType := reflect.FuncOf(inTypes, outTypes, false)
	fn := reflect.MakeFunc(newFuncType, func(in []reflect.Value) []reflect.Value {
		results := b.V.Call(in)
		return reflect.ValueOf(fmap).Call(results)
	})
	return newFunctionBoxed(fn.Interface())
}

func (b *_FunctionBoxed) FlatMap(fmap interface{}) Boxed {
	return b
}

func (b *_FunctionBoxed) Unbox(dest interface{}) {
	value := reflect.ValueOf(dest)
	panicCondition(value.Kind() != reflect.Ptr, "gomads: Unbox needs a pointer")
	panicCondition(value.IsNil(), "gomads: Unbox (nil)")

	fn, err := baseType(value.Type(), reflect.Func)
	panicCondition(err != nil, "gomads: Unbox (not function destination)")
	panicCondition(fn.NumIn() != b.T.NumIn(), "gomads: Unbox (not equal input params)")
	panicCondition(fn.NumOut() != b.T.NumOut(), "gomads: Unbox (not equal input params)")
	for i := 0; i < fn.NumIn(); i++ {
		panicCondition(fn.In(i) != b.T.In(i), "gomads: Unbox (not same function signature)")
	}
	for i := 0; i < fn.NumOut(); i++ {
		panicCondition(fn.Out(i) != b.T.Out(i), "gomads: Unbox (not same function signature)")
	}

	reflect.Indirect(value).Set(b.V)
}
